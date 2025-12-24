# Deadlock Detection System

This project implements a distributed deadlock detection system using the **Chandy-Lamport Snapshot Algorithm**. It simulates a set of concurrent processes managing resources and detects if they have entered a deadlocked state.

## How it Works

The system consists of three main components:

1.  **Processes**: Independent entities that own resources and can request resources from others. They communicate via message passing (Go channels).
2.  **Snapshot Algorithm (Chandy-Lamport)**: A mechanism to capture a consistent global state of the distributed system (processes + communication channels) without freezing the application.
3.  **Snapshot Server (The Detector)**: A centralized observer that collects local snapshots from all processes, constructs a **Wait-For Graph (WFG)**, and analyzes it for cycles.

## Simulation Logic (`main.go`)

The `main.go` file orchestrates a classic circular deadlock scenario to demonstrate the detection capability:

1.  **Initialization**:
    *   A `SnapshotServer` is started to listen for reports.
    *   Three processes are created:
        *   **Process A** (owns Resource 1)
        *   **Process B** (owns Resource 2)
        *   **Process C** (owns Resource 3)
    *   Processes are fully connected with bidirectional channels.

2.  **Deadlock Creation**:
    *   **A** requests Resource 2 (held by **B**).
    *   **B** requests Resource 3 (held by **C**).
    *   **C** requests Resource 1 (held by **A**).
    *   This establishes a circular dependency: `A -> B -> C -> A`. Since no process releases its resource, the system hangs.

3.  **Detection**:
    *   **Process A** initiates the Chandy-Lamport snapshot by recording its state and flooding `MARKER` messages.
    *   As markers propagate, every process records its local state (who it is waiting for) and the state of its incoming channels (messages in transit).
    *   Once a process has received markers on all input links, it sends a `SnapshotReport` to the central server.

4.  **Analysis**:
    *   The **Snapshot Server** aggregates all reports.
    *   It builds a global **Wait-For Graph (WFG)**. Edges represent dependencies found in:
        *   **Process State**: Explicitly waiting for a grant.
        *   **Channel State**: Request messages that were "on the wire" during the snapshot.
    *   A Depth-First Search (DFS) runs on the WFG. If a cycle is found, it confirms the deadlock.

## Build and Run

```bash
go build .
./deadlock-detection-system
```

### Expected Output

```text
...
[Server] Received snapshot from A
[Server] Received snapshot from B
[Server] Received snapshot from C
[Server] All reports received. Building Wait-For Graph...
[WFG] Edge A -> B (Process State)
[WFG] Edge B -> C (Process State)
[WFG] Edge C -> A (Process State)
DEADLOCK DETECTED: Cycle found: [A B C] -> A
```