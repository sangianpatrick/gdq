package goredisv8

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sangianpatrick/gdq"
)

type sender struct {
	rc redis.UniversalClient
}

func NewSender(rc redis.UniversalClient) (s gdq.Sender, err error) {
	if rc == nil {
		return nil, gdq.ErrNoRedisClient
	}

	s = &sender{
		rc: rc,
	}

	return
}

func (s *sender) Send(ctx context.Context, topic string, delay time.Duration, message *gdq.Message) (err error) {
	delayTime := time.Now().Add(delay)
	delayTimeMS := delayTime.UnixNano() / int64(time.Millisecond)

	message.Score = float64(delayTimeMS)
	messageByte, err := json.Marshal(message)
	if err != nil {
		err = gdq.ErrInvalidMessageToSerialize
		return
	}

	z := &redis.Z{
		Score:  float64(delayTimeMS),
		Member: messageByte,
	}

	_, err = s.rc.ZAdd(ctx, topic, z).Result()
	if err != nil {
		err = gdq.ErrRedisCommandError
		return
	}

	return
}
