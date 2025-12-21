package main

import "fmt"

type MessageType string

const (
	RequestLock MessageType = "REQUEST_LOCK"
	GrantLock   MessageType = "GRANT_LOCK"
	Marker      MessageType = "MARKER"
)

type Message struct {
	From       string
	To         string
	Type       MessageType
	ResourceID int
}

type Process struct {
	ID        string
	Resources map[int]bool
	Inbound   map[string]<-chan Message
	Outbound  map[string]chan<- Message
}

func NewProcess(id string, resources []int) *Process {
	owned := make(map[int]bool, len(resources))
	for _, r := range resources {
		owned[r] = true
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

func main() {
	a := NewProcess("A", []int{1})
	b := NewProcess("B", []int{2})
	c := NewProcess("C", []int{3})

	processes := []*Process{a, b, c}
	ConnectFully(processes)

	fmt.Printf("initialized %d processes\n", len(processes))
}
