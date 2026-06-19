package kafkaClient

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	w *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) KafkaProducer {
	return KafkaProducer{w: &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}}
}

func (p KafkaProducer) Send(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	return p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: data,
	})
}

func (p KafkaProducer) Close() error {
	return p.w.Close()
}
