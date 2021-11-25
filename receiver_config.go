package gdq

import "time"

// ReceiverConfig is a configuration
type ReceiverConfig struct {
	// Handler which handles the delayed message
	Handler Handler
	// MaxBatch is a total message that will be process at a time, the default is 5
	MaxBatch int
	// FetchInterval is time interval to fetch message (millisecond), the default is 200
	FetchInterval int
	// FetchTimeout is max wait time before timeout error (milliesecond), the default is 2000
	FetchTimeout int
	// CloseTime is time that will be spent while closing the subscriber (millisecond), the default is 500
	CloseTime int
	// FetchBeforTime defines the consumer to polls the message which scheduled before or exact the given time.
	// Leave it empty to fetch before now (current time)
	FetchBeforeTime *time.Time
}
