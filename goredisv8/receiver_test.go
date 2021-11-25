package goredisv8_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/sangianpatrick/gdq"
	"github.com/sangianpatrick/gdq/goredisv8"
)

type mockHandler struct{}

func (mockHandler) Handle(ctx context.Context, message interface{}) (err error) {
	return nil
}

var handler = new(mockHandler)

var topic = "test-topic"

func TestReceiver_Success(t *testing.T) {
	message := "hello"
	fetchBeforeTime := time.Now().Add(time.Hour * 1)
	maxBatch := 5
	rc, mock := redismock.NewClientMock()

	max := fetchBeforeTime.UnixNano() / 1_000_000

	mockZSlice := mock.ExpectZRangeByScoreWithScores(topic, &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprint(max),
		Offset: 0,
		Count:  int64(maxBatch),
	})

	mockZSlice.SetVal([]redis.Z{
		{
			Score:  float64(max),
			Member: message,
		},
	})
	mockZSlice.SetErr(nil)

	mockInt := mock.ExpectZRem(topic, message)
	mockInt.SetVal(1)

	rec, _ := goredisv8.NewReceiver(rc, gdq.ReceiverConfig{
		Handler:         handler,
		MaxBatch:        maxBatch,
		FetchBeforeTime: &fetchBeforeTime,
	}, topic)

	rec.Receive()

	time.Sleep(time.Second * 2)

	rec.Close()

	mock.ClearExpect()
}
