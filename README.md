
# Redis Pub/Sub At Least Once Delivery Example

this example demonstrate how to make Redis that are by default only have two mode:
- broadcast mode (at most once delivery, so if client down at the moment of broadcast time they will not ever receive the message)
- queue mode (at least once delivery, but not broadcasting/without consumer group, so 1 message only acked by 1 subscriber)

to mimic Kafka/Jetsream/standard-AMQP, we need to have both: broadcast (having a customer group) and at least once delivery (ack per consumer group)

## How to run

```shell
docker compose up

make run CMD='go run main.go publisher topic1'
```

on another terminal:

```shell
make run CMD='go run main.go subscriber topic1'

make run CMD='go run main.go subscriber topic1 subscriber1'
```

## Maintenance checklist

- [x] Go runtime updated to 1.26.5.
- [x] Redis Pub/Sub, Go Cloud, and gotro dependencies refreshed.
- [x] `make test` compiles the example without requiring Redis.
- [x] `make verify-dependency-security` and `make vulncheck` check dependency security.
