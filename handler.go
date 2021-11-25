package gdq

import "context"

// Handler handles the stream of delayed message
type Handler interface {
	// Handle will run the task after consume the delayed message and process the message
	Handle(ctx context.Context, message interface{}) (err error)
}
