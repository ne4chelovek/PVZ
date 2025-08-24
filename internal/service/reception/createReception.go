package reception

import (
	"PVZ/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (s *ReceptionService) CreateReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	rec, err := s.Repo.GetOpenReceptionForPVZ(ctx, pvzID)
	if err != nil {
		return nil, err
	}
	if rec != nil && rec.Status == "in_progress" {
		return nil, fmt.Errorf("reception already exists")
	}

	reception := &model.Reception{
		ID:       uuid.New().String(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "in_progress",
	}

	return s.Repo.Create(ctx, reception)
}
