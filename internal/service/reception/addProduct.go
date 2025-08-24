package reception

import (
	"PVZ/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (s *ReceptionService) AddProduct(ctx context.Context, pvzID, productType string) (*model.Product, error) {
	validTypes := map[string]bool{"электроника": true, "одежда": true, "обувь": true}
	if !validTypes[productType] {
		return nil, fmt.Errorf("неверный тип товара")
	}

	rec, err := s.Repo.GetOpenReceptionForPVZ(ctx, pvzID)
	if err != nil || rec.Status != "in_progress" {
		return nil, fmt.Errorf("нет активной приёмки")
	}

	product := &model.Product{
		ID:          uuid.New().String(),
		DateTime:    time.Now(),
		Type:        productType,
		ReceptionID: rec.ID,
	}

	if err := s.Repo.AddProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}
