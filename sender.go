package gdq

import (
	"context"
	"time"
)

// Sender is a collection of behavior of sender.
type Sender interface {
	// Send will put the message into queue depends on its delayed time.
	Send(ctx context.Context, topic string, delay time.Duration, message *Message) (err error)
}
