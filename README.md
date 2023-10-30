
# Redis Pub/Sub At Least Once Delivery Example

this example demo-ing how to make Redis That by default only have two mode:
- broadcast mode (at most once delivery)
- queue mode (at last once delivery)

to mimic Kafka/Jetsream/standard-AMQP, we need to have both: broadcast (having a customer group) and at last once delivery (ack per consumer group)

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