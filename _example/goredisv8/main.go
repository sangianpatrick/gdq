package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sangianpatrick/gdq"
	"github.com/sangianpatrick/gdq/goredisv8"
)

type myHandler struct{}

func (*myHandler) Handle(ctx context.Context, message interface{}) (err error) {
	m, ok := message.(*gdq.Message)
	if !ok {
		return fmt.Errorf("unprocessable entity")
	}

	for headerKey, headerValues := range m.Headers {
		for _, headerValue := range headerValues {
			fmt.Printf("[Message Header] Key: %s | Value: %s\n", headerKey, headerValue)
		}
	}

	fmt.Printf("[Message Key] %s\n", m.Key)
	fmt.Printf("[Message Value] %s\n", string(m.Value))
	fmt.Printf("[Message Score] %f\n", m.Score)
	fmt.Println()

	return
}

var topic string = "example-topic"

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := rc.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	receiver, err := goredisv8.NewReceiver(rc, gdq.ReceiverConfig{
		Handler: new(myHandler),
	}, topic)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		receiver.Receive()
	}()

	sender, err := goredisv8.NewSender(rc)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			var delay time.Duration
			msg := fmt.Sprintf("Message - %d", i)
			schedule := time.Date(2021, time.December, 3, 2, 56, 0, 0, time.Local).UnixNano() / int64(time.Millisecond)
			delay = time.Duration(schedule-(time.Now().UnixNano()/int64(time.Millisecond))) * time.Millisecond

			sender.Send(context.Background(), topic, delay, &gdq.Message{
				Key: fmt.Sprintf("key:%d", i),
				Headers: map[string][]string{
					"number": []string{fmt.Sprint(i)},
				},
				Value: []byte(msg),
			})
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	receiver.Close()
	rc.Close()
}
