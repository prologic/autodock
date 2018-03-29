package events

import (
	etypes "github.com/docker/docker/api/types/events"
)

// Handler ...
type Handler interface {
	Handle(message *Message) error
}

// Message ...
type Message struct {
	etypes.Message
}
