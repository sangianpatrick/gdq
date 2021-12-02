# [gdq] Golang Delay Queue
GDQ is a library that leverage db or cache to be setup as a delay queue. For current version, Only redis can adapt to this library.

# Example
## Receiver
### go-redis v8
```go
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

sigterm := make(chan os.Signal, 1)
signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
<-sigterm

receiver.Close()
rc.Close()
```
