package outbox

import (
	"PVZ/internal/kafka"
	"PVZ/internal/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxProcessor struct {
	eventRepo    repository.EventRepository
	kafkaProd    *kafka.KafkaProducer
	db           *pgxpool.Pool
	pollInterval time.Duration
	batchSize    int
}

func NewOutboxProcessor(eventRepo repository.EventRepository, kafkaProd *kafka.KafkaProducer, db *pgxpool.Pool) *OutboxProcessor {
	return &OutboxProcessor{
		eventRepo:    eventRepo,
		kafkaProd:    kafkaProd,
		db:           db,
		pollInterval: 5 * time.Second,
		batchSize:    100,
	}
}
