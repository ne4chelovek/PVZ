package kafka

import (
	"PVZ/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
	"time"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (k *KafkaProducer) PublishEvent(ctx context.Context, eventType string, aggregateID string, payload []byte) error {
	event := model.Event{
		ID:          uuid.New().String(),
		EventType:   eventType,
		AggregateId: aggregateID,
		Payload:     string(payload),
		CreatedAt:   time.Now(),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(aggregateID),
		Value: sarama.StringEncoder(eventBytes),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf("Event published to Kafka - Topic: %s, Partition: %d, Offset: %d", k.topic, partition, offset)
	return nil
}

func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}
