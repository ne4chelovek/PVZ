package event

import (
	"PVZ/internal/model"
	"PVZ/internal/repository"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type eventRepository struct {
	db repository.QueryRunner
}

func NewEventRepository(db *pgxpool.Pool) repository.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) WithTx(tx pgx.Tx) repository.EventRepository {
	return &eventRepository{db: tx}
}

func (r *eventRepository) Event(ctx context.Context, event *model.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	query := squirrel.Insert("outbox_events").
		Columns("id", "aggregate_id", "event_type", "payload", "created_at").
		Values(event.ID, event.AggregateId, event.EventType, event.Payload, event.CreatedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *eventRepository) GetUnprocessedEvents(ctx context.Context, tx pgx.Tx, limit int) ([]*model.Event, error) {
	query := squirrel.Select("id", "aggregate_id", "event_type", "payload", "created_at").
		From("outbox_events").
		Where(squirrel.Eq{"processed": false}).
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		var event model.Event
		if err := rows.Scan(&event.ID, &event.AggregateId, &event.EventType, &event.Payload, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) MarkAsProcessed(ctx context.Context, tx pgx.Tx, eventID string) error {
	query := squirrel.Update("outbox_events").
		Set("processed", true).
		Where(squirrel.Eq{"id": eventID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	return err
}
