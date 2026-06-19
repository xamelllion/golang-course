package kafkaClient

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	r *kafka.Reader
}

func NewKafkaReader(brokers []string, topic, groupID string) KafkaConsumer {
	return KafkaConsumer{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (p KafkaConsumer) Run(ctx context.Context, handle func(context.Context, []byte) error) error {
	for {
		msg, err := p.r.FetchMessage(ctx)
		if err != nil {
			return err
		}

		if err := handle(ctx, msg.Value); err != nil {
			log.Printf("error in handle %v", err)
		}

		if err := p.r.CommitMessages(ctx, msg); err != nil {
			log.Printf("error in commit message %v", err)
			return err
		}
	}
}

func (p KafkaConsumer) Close() error {
	return p.r.Close()
}
