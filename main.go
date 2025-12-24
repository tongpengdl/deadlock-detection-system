package main

import (
	"fmt"
	"time"
)

func main() {
	// Initialize Snapshot Server
	// Expecting 3 reports
	server := NewSnapshotServer(3)

	a := NewProcess("A", []int{1}, server)
	b := NewProcess("B", []int{2}, server)
	c := NewProcess("C", []int{3}, server)

	processes := []*Process{a, b, c}
	ConnectFully(processes)

	fmt.Printf("initialized %d processes\n", len(processes))

	for _, process := range processes {
		process.Start()
	}

	// Create a deadlock scenario
	// A holds 1, wants 2 (held by B)
	// B holds 2, wants 3 (held by C)
	// C holds 3, wants 1 (held by A)
	a.RequestLock(2, "B")
	b.RequestLock(3, "C")
	c.RequestLock(1, "A")

	time.Sleep(500 * time.Millisecond)
	fmt.Println("deadlock likely: each process holds one resource and waits on another")

	fmt.Println("--- Initiating Snapshot from Process A ---")
	a.RecordState()

	// Allow time for snapshot messages to propagate and server to process
	time.Sleep(2 * time.Second)
}
