package reception

import (
	"PVZ/internal/model"
	"context"
	"fmt"
)

func (s *ReceptionService) CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	rec, err := s.Repo.GetOpenReceptionForPVZ(ctx, pvzID)
	if err != nil || rec.Status != "in_progress" {
		return nil, fmt.Errorf("нет активной приёмки для закрытия")
	}

	rec.Status = "close"
	if err := s.Repo.CloseReception(ctx, rec.ID); err != nil {
		return nil, err
	}

	return rec, nil
}
