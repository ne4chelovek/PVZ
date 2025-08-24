package reception

import (
	"context"
	"fmt"
)

func (s *ReceptionService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	rec, err := s.Repo.GetOpenReceptionForPVZ(ctx, pvzID)
	if err != nil {
		return fmt.Errorf("нет активной приёмки: %w", err)
	}
	if rec == nil || rec.Status != "in_progress" {
		return fmt.Errorf("нет активной приёмки")
	}

	products, err := s.Repo.GetProductsByReception(ctx, rec.ID)
	if err != nil {
		return fmt.Errorf("ошибка при получении товаров: %w", err)
	}
	if len(products) == 0 {
		return fmt.Errorf("нет товаров для удаления")
	}

	last := products[len(products)-1]
	return s.Repo.DeleteLastProductByPVZ(ctx, last.ID)
}
