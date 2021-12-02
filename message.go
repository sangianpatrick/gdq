package gdq

// Message is an object that can be sent and received.
type Message struct {
	// Key is a key of a message. Make sure it is globally unique to avoid deleting queue with the same message instance.
	Key string `json:"key"`
	// Headers will describe the information of a message. Make sure to fill it.
	Headers map[string][]string `json:"headers"`
	// Value is a value of a message.
	Value []byte `json:"value"`
	// Score is a score of message that used for sorting the message.
	Score float64 `json:"score"`
}
