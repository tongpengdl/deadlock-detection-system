package main

import (
	"fmt"
	"sort"
	"sync"
)

type SnapshotServer struct {
	mu            sync.Mutex
	expectedCount int
	reports       map[string]SnapshotReport
}

func NewSnapshotServer(expectedCount int) *SnapshotServer {
	return &SnapshotServer{
		expectedCount: expectedCount,
		reports:       make(map[string]SnapshotReport),
	}
}

// Collect receives a snapshot report from a process.
func (s *SnapshotServer) Collect(report SnapshotReport) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("[Server] Received snapshot from %s\n", report.ProcessID)
	s.reports[report.ProcessID] = report

	if len(s.reports) == s.expectedCount {
		fmt.Println("[Server] All reports received. Building Wait-For Graph...")
		s.DetectDeadlock()
	}
}

// DetectDeadlock constructs the WFG and checks for cycles.
func (s *SnapshotServer) DetectDeadlock() {
	// 1. Build Adjacency List (WFG)
	// Map: ProcessID -> List of ProcessIDs it is waiting for
	adj := make(map[string][]string)

	// Initialize nodes
	for pid := range s.reports {
		adj[pid] = []string{}
	}

	for pid, report := range s.reports {
		// A waiting for B (Local State)
		for target, isWaiting := range report.RecordedWaitingFor {
			if isWaiting {
				fmt.Printf("[WFG] Edge %s -> %s (Process State)\n", pid, target)
				adj[pid] = append(adj[pid], target)
			}
		}

		// A waiting for B (Channel State)
		// Channel messages are incoming messages to 'pid'.
		// If 'pid' recorded a RequestLock from 'sender', it means 'sender' sent a request
		// that was in-flight. Effectively, 'sender' is waiting for 'pid'.
		// Edge: sender -> pid
		for sender, msgs := range report.RecordedChannelMessages {
			for _, msg := range msgs {
				if msg.Type == RequestLock {
					fmt.Printf("[WFG] Edge %s -> %s (Channel State: Request in-flight)\n", sender, pid)
					adj[sender] = append(adj[sender], pid)
				}
			}
		}
	}

	// 2. Cycle Detection (DFS)
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	hasDeadlock := false

	// Helper function for DFS
	var dfs func(u string, path []string)
	dfs = func(u string, path []string) {
		visited[u] = true
		recStack[u] = true
		path = append(path, u)

		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v, path)
			} else if recStack[v] {
				// Cycle detected
				hasDeadlock = true
				fmt.Printf("DEADLOCK DETECTED: Cycle found: %v -> %s\n", path, v)
			}
		}
		recStack[u] = false
	}

	// Keys should be sorted for deterministic iteration order (good for testing/logs)
	var keys []string
	for k := range adj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, u := range keys {
		if !visited[u] {
			dfs(u, []string{})
		}
	}

	if !hasDeadlock {
		fmt.Println("No deadlock detected.")
	}
}
