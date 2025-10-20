package outbox

import (
	"PVZ/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func (p *OutboxProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.processBatch(ctx); err != nil {
				log.Printf("Error processing outbox batch: %v", err)
			}
		}
	}
}

func (p *OutboxProcessor) processBatch(ctx context.Context) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	events, err := p.getUnprocessedEvents(ctx, tx)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := p.processEvent(ctx, tx, event); err != nil {
			log.Printf("Failed to process event %s: %v", event.ID, err)
			continue
		}
	}

	return tx.Commit(ctx)
}

func (p *OutboxProcessor) getUnprocessedEvents(ctx context.Context, tx pgx.Tx) ([]*model.Event, error) {
	return p.eventRepo.GetUnprocessedEvents(ctx, tx, p.batchSize)
}

func (p *OutboxProcessor) processEvent(ctx context.Context, tx pgx.Tx, event *model.Event) error {
	if err := p.kafkaProd.PublishEvent(ctx, event.EventType, event.AggregateId, event.Payload); err != nil {
		return fmt.Errorf("failed to publish to kafka: %w", err)
	}
	return p.eventRepo.MarkAsProcessed(ctx, tx, event.ID)
}
