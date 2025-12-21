package main

import "fmt"

type ResourceLock struct {
	ID     int
	Holder string
	Queue  []string
}

type Process struct {
	ID        string
	Resources map[int]*ResourceLock
	Inbound   map[string]<-chan Message
	Outbound  map[string]chan<- Message
	inbox     chan Message
}

func NewProcess(id string, resources []int) *Process {
	owned := make(map[int]*ResourceLock, len(resources))
	for _, r := range resources {
		owned[r] = &ResourceLock{
			ID:     r,
			Holder: id,
		}
	}

	return &Process{
		ID:        id,
		Resources: owned,
		Inbound:   make(map[string]<-chan Message),
		Outbound:  make(map[string]chan<- Message),
	}
}

const defaultChannelBuffer = 100

func connectDirected(from, to *Process, ch chan Message) {
	from.Outbound[to.ID] = ch
	to.Inbound[from.ID] = ch
}

func ConnectBidirectional(a, b *Process) {
	chAB := make(chan Message, defaultChannelBuffer)
	chBA := make(chan Message, defaultChannelBuffer)
	connectDirected(a, b, chAB)
	connectDirected(b, a, chBA)
}

func ConnectFully(processes []*Process) {
	for i := 0; i < len(processes); i++ {
		for j := i + 1; j < len(processes); j++ {
			ConnectBidirectional(processes[i], processes[j])
		}
	}
}

func (p *Process) Start() {
	p.inbox = make(chan Message, defaultChannelBuffer)
	for _, inbound := range p.Inbound {
		ch := inbound
		go func() {
			for msg := range ch {
				p.inbox <- msg
			}
		}()
	}

	go func() {
		for msg := range p.inbox {
			p.handleMessage(msg)
		}
	}()
}

func (p *Process) RequestLock(resourceID int, ownerID string) {
	ch, ok := p.Outbound[ownerID]
	if !ok {
		fmt.Printf("process %s has no link to %s for resource %d\n", p.ID, ownerID, resourceID)
		return
	}

	msg := Message{
		From:       p.ID,
		To:         ownerID,
		Type:       RequestLock,
		ResourceID: resourceID,
	}
	ch <- msg
	fmt.Printf("process %s requesting resource %d from %s\n", p.ID, resourceID, ownerID)
}

func (p *Process) handleMessage(msg Message) {
	switch msg.Type {
	case RequestLock:
		p.handleRequestLock(msg)
	case GrantLock:
		p.handleGrantLock(msg)
	case Marker:
		fmt.Printf("process %s received marker from %s\n", p.ID, msg.From)
	default:
		fmt.Printf("process %s received unknown message type %q from %s\n", p.ID, msg.Type, msg.From)
	}
}

func (p *Process) handleRequestLock(msg Message) {
	lock, ok := p.Resources[msg.ResourceID]
	if !ok {
		fmt.Printf("process %s cannot grant resource %d to %s: not owner\n", p.ID, msg.ResourceID, msg.From)
		return
	}

	if lock.Holder == "" {
		lock.Holder = msg.From
		p.sendGrant(msg.From, msg.ResourceID)
		fmt.Printf("process %s granted resource %d to %s\n", p.ID, msg.ResourceID, msg.From)
		return
	}

	lock.Queue = append(lock.Queue, msg.From)
	fmt.Printf(
		"process %s queued request from %s for resource %d (held by %s)\n",
		p.ID,
		msg.From,
		msg.ResourceID,
		lock.Holder,
	)
}

func (p *Process) handleGrantLock(msg Message) {
	fmt.Printf("process %s received grant for resource %d from %s\n", p.ID, msg.ResourceID, msg.From)
}

func (p *Process) sendGrant(to string, resourceID int) {
	ch, ok := p.Outbound[to]
	if !ok {
		fmt.Printf("process %s cannot send grant for resource %d to %s: no link\n", p.ID, resourceID, to)
		return
	}

	ch <- Message{
		From:       p.ID,
		To:         to,
		Type:       GrantLock,
		ResourceID: resourceID,
	}
}
