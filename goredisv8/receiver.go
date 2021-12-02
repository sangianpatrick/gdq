package goredisv8

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sangianpatrick/gdq"
)

type receiver struct {
	rc              redis.UniversalClient
	handler         gdq.Handler
	closed          chan struct{}
	topic           string
	fetchInterval   time.Duration
	fetchTimeout    time.Duration
	closeTime       time.Duration
	fetchBeforeTime *time.Time
}

// NewReceiver is a constructor. Make sure the receiver only instantiate once globally accross services for a topic.
func NewReceiver(client redis.UniversalClient, config gdq.ReceiverConfig, topic string) (rec gdq.Receiver, err error) {
	closed := make(chan struct{}, 1)
	defaultFetchInterval := 200
	defaultFetchTimeout := 2000
	defaultCloseTime := 500

	if client == nil {
		return nil, gdq.ErrNoRedisClient
	}

	if config.Handler == nil {
		return nil, gdq.ErrNoHandler
	}

	if config.FetchInterval <= 0 {
		config.FetchInterval = defaultFetchInterval
	}

	if config.FetchTimeout <= 0 {
		config.FetchTimeout = defaultFetchTimeout
	}

	if config.CloseTime <= 0 {
		config.CloseTime = defaultCloseTime
	}

	rec = &receiver{
		rc:              client,
		handler:         config.Handler,
		closed:          closed,
		topic:           topic,
		fetchInterval:   time.Duration(defaultFetchInterval) * time.Millisecond,
		fetchTimeout:    time.Duration(defaultFetchTimeout) * time.Millisecond,
		closeTime:       time.Duration(defaultCloseTime) * time.Millisecond,
		fetchBeforeTime: config.FetchBeforeTime,
	}

	return
}

func (r *receiver) Receive() {
	go func() {
		for {
			select {
			case <-r.closed:
				return
			case <-time.After(r.fetchInterval):
				ctx, cancel := context.WithTimeout(context.Background(), r.fetchTimeout)

				bunchOfZs, err := r.pool(ctx)
				if err != nil {
					cancel()
					log.Println(err)
					continue
				}

				if len(bunchOfZs) > 0 {
					lastZ := bunchOfZs[len(bunchOfZs)-1]
					err := r.deleteZ(ctx, lastZ)
					if err != nil {
						cancel()
						log.Println(err)
						continue
					}
				}

				cancel()

				for _, z := range bunchOfZs {
					handlerContext := context.Background()

					msgRaw, ok := z.Member.(string)

					if !ok {
						fmt.Println("invalid message")
						continue
					}

					m := new(gdq.Message)
					json.Unmarshal([]byte(msgRaw), m)
					r.handler.Handle(handlerContext, m)
				}

			}
		}
	}()
}

func (r *receiver) Close() {
	close(r.closed)
	<-time.After(r.closeTime)
}

func (r *receiver) pool(ctx context.Context) (bunchOfZs []redis.Z, err error) {
	max := time.Now().UnixNano() / 1_000_000

	if r.fetchBeforeTime != nil {
		max = r.fetchBeforeTime.UnixNano() / 1_000_000
	}

	zSliceCmd := r.rc.ZRangeByScoreWithScores(ctx, r.topic, &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprint(max),
		Offset: 0,
		Count:  1,
	})

	bunchOfZs, err = zSliceCmd.Result()
	return
}

func (r *receiver) deleteZ(ctx context.Context, bunchOfZs ...redis.Z) (err error) {

	// use this approach to delete by score range
	//
	// intScore := int64(z.Score)
	// strScore := fmt.Sprintf("%d", intScore)
	// fmt.Println(strScore)
	// intCmd := c.rc.ZRemRangeByScore(context.Background(), c.topic, strScore, strScore)
	// _, err = intCmd.Result()

	// use this approach to delete by member
	zMember := make([]interface{}, len(bunchOfZs))
	for i, z := range bunchOfZs {
		zMember[i] = z.Member
	}

	intCmd := r.rc.ZRem(ctx, r.topic, zMember...)
	_, err = intCmd.Result()

	return
}
