package gdq

import "fmt"

var (
	ErrNoHandler                 = fmt.Errorf("gdq: no handler specified")
	ErrNoRedisClient             = fmt.Errorf("gdq: no redis client specified")
	ErrRedisCommandError         = fmt.Errorf("gdq: redis error")
	ErrInvalidMessageToSerialize = fmt.Errorf("gdq: invalid message to serialize")
)
