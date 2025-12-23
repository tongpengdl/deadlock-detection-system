package main

import (
	"fmt"
	"time"
)

func main() {
	a := NewProcess("A", []int{1})
	b := NewProcess("B", []int{2})
	c := NewProcess("C", []int{3})

	processes := []*Process{a, b, c}
	ConnectFully(processes)

	fmt.Printf("initialized %d processes\n", len(processes))

	for _, process := range processes {
		process.Start()
	}

	a.RequestLock(2, "B")
	b.RequestLock(3, "C")
	c.RequestLock(1, "A")

	time.Sleep(500 * time.Millisecond)
	fmt.Println("deadlock likely: each process holds one resource and waits on another")

	fmt.Println("--- Initiating Snapshot from Process A ---")
	a.RecordState()

	// Allow time for snapshot messages to propagate
	time.Sleep(1 * time.Second)
}
