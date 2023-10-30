package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/covrom/redispubsub"
	"github.com/kokizzu/gotro/L"
	"gocloud.dev/pubsub"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println(`required argument: publisher/subscriber topicPrefix [subscriberName or random]`)
		fmt.Println(`example: publisher topic1`)
		fmt.Println(`example: subscriber topic1 subscriber1`) // subscriber = consumer in this case since we're not using worker group
		return
	}

	mode := os.Args[1]
	topic := os.Args[2]

	switch mode {
	case `publisher`:
		err := os.Setenv(`REDIS_URL`, `redis://localhost:6379`)
		L.PanicIf(err, `os.Setenv REDIS_URL`)
		ctx := context.Background()
		pub, err := pubsub.OpenTopic(ctx, fmt.Sprintf("redis://topics/%s", topic))
		defer pub.Shutdown(ctx)
		L.PanicIf(err, `pubsub.OpenTopic`)
		counter := uint64(0)
		for {
			msg := fmt.Sprintf(`hello %d`, counter)
			atomic.AddUint64(&counter, 1)
			fmt.Printf("publishing to topic: %s %s\n", topic, msg)
			m := &pubsub.Message{
				Body: []byte(msg),
				Metadata: map[string]string{
					`createdAt`: time.Now().Format(time.DateTime),
				},
			}

			err = pub.Send(ctx, m)
			L.IsError(err, `topic.Send`)
			time.Sleep(1 * time.Second)
		}
	case `subscriber`:
		var subscriber string
		if len(os.Args) <= 3 {
			// random subscriber name
			subscriber = fmt.Sprintf(`sub%d`, time.Now().UnixNano()%1000)
		} else {
			subscriber = os.Args[3]
		}
		// connect to replica
		// cannot write (Ack) to readonly replica
		//port := 6380 + rand.Intn(2) // connect to random redis replica
		//hostPort := fmt.Sprintf(`redis://localhost:%d`, port)
		// so we have to connect to master
		hostPort := `redis://localhost:6379`
		err := os.Setenv(`REDIS_URL`, hostPort)
		L.PanicIf(err, `os.Setenv REDIS_URL`)
		ctx := context.Background()
		// first %s is group, but since we only have 1 subscriber per 1 group, we use group=subscsriber name
		subs, err := pubsub.OpenSubscription(ctx, fmt.Sprintf("redis://%s?consumer=%s&topic=topics/%s", subscriber, subscriber, topic))
		L.PanicIf(err, `pubsub.OpenSubscription`)
		defer subs.Shutdown(ctx)

		fmt.Printf("%s subscribing to topic %s\n", subscriber, topic)
		for {
			msg, err := subs.Receive(ctx)
			if err != nil {
				fmt.Printf("%s got error: %s\n", subscriber, err)
			} else {
				fmt.Printf("%s got message: %q\n", subscriber, msg.Body)
				msg.Ack()
			}
		}
	}
}
