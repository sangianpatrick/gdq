package gdq

// Message is an object that can be sent and received.
type Message struct {
	Key     string              `json:"key"`
	Headers map[string][]string `json:"headers"`
	Value   []byte              `json:"value"`
	Score   float64             `json:"score"`
}
