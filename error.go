package gdq

import "fmt"

var (
	ErrNoHandler     = fmt.Errorf("gdq: no handler specified")
	ErrNoRedisClient = fmt.Errorf("gdq: no redis client specified")
)
