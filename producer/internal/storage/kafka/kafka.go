package kafka

import (
	"github.com/segmentio/kafka-go"
)

type Storage struct {
	writer *kafka.Writer
}

func New(addr, topic string) *Storage {
	const op = "storage.kafka.New"

	producer := &kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Storage{writer: producer}
}
