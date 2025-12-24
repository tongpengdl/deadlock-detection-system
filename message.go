package main

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

// SnapshotReport contains the local state of a process after a snapshot.
type SnapshotReport struct {
	ProcessID               string
	RecordedResources       map[int]*ResourceLock
	RecordedWaitingFor      map[string]bool
	RecordedChannelMessages map[string][]Message
}
