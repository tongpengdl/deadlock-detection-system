# deadlock-detection-system

Prototype of a deadlock detection system inspired by Chandy-Lamport snapshots.

## Build and run
```bash
go build .
./deadlock-detection-system
```

## Milestone 1.1: Process & messages
- `Process` has an ID, owned resources, and per-link inbound/outbound channels.
- Messages carry a type (`REQUEST_LOCK`, `GRANT_LOCK`, `MARKER`) and resource ID.
- Each directed link has its own buffered Go channel (no global channel).

## Milestone 1.2: Grasping logic (next)
- Processes request resources they do not own.
- Owners grant if free; otherwise they queue the request.
- Run three processes (A, B, C) with resources 1, 2, 3 to observe deadlock.
