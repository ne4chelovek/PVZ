package pvz

import (
	"PVZ/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

func (s *PVZService) CreatePVZ(ctx context.Context, pvz *model.PVZ) (*model.PVZ, error) {
	if pvz.ID == "" {
		pvz.ID = uuid.New().String()
	}

	if pvz.RegistrationDate.IsZero() {
		pvz.RegistrationDate = time.Now()
	}

	validCities := map[string]bool{
		"Москва":          true,
		"Санкт-Петербург": true,
		"Казань":          true,
	}
	if !validCities[pvz.City] {
		return nil, fmt.Errorf("city must be one of: Москва, Санкт-Петербург, Казань")
	}

	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txPVZRepo := s.pvzRepo.WithTx(tx)
	txEventRepo := s.outBox.WithTx(tx)

	if err := txPVZRepo.Create(ctx, pvz); err != nil {
		return nil, fmt.Errorf("failed to create PVZ in transaction: %w", err)
	}
	event := &model.Event{
		AggregateId: pvz.ID,
		EventType:   "PVZCreated",
		Payload:     `{"id": "` + pvz.ID + `", "city": "` + pvz.City + `"}`, // или через json.Marshal
	}

	if err := txEventRepo.Event(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to save event to outbox: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return pvz, nil
}
