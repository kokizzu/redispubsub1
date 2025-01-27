
# Redis Pub/Sub At Least Once Delivery Example

this example demo-ing how to make Redis that are by default only have two mode:
- broadcast mode (at most once delivery, so if client down at the moment of broadcast time they will not ever receive the message)
- queue mode (at last once delivery, but not broadcasting/without consumer group, so 1 message only acked by 1 subscriber)

to mimic Kafka/Jetsream/standard-AMQP, we need to have both: broadcast (having a customer group) and at least once delivery (ack per consumer group)

## How to run

```shell
docker compose up

go run main.go publisher topic1
```

on another terminal:

```shell
go run main.go subscriber topic1

go run main.go subscriber topic1 subscriber1
```
