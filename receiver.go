package gdq

// Receiver is a collection of behavior of receiver.
type Receiver interface {
	// Receive will pull the queue from stream.
	// It already run with goroutine
	Receive()
	// Close will close the instance
	Close()
}
