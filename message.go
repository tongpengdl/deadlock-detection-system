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
