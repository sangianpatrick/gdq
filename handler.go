package gdq

import "context"

// Handler handles the stream of delayed message
type Handler interface {
	// Handle will run the task after consume the delayed message and process the message.
	// Make sure the implementation type assert the message as *gdq.Message.
	Handle(ctx context.Context, message interface{}) (err error)
}
