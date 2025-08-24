package pvz

import (
	"PVZ/internal/model"
	"context"
	"fmt"
)

func (s *PVZService) GetAllPVZWithReceptions(ctx context.Context, startDate, endDate *string, page, limit int) ([]*model.PVZWithReceptions, error) {
	pvzs, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get PVZ: %w", err)
	}

	var result []*model.PVZWithReceptions
	for _, p := range pvzs {
		receptions, err := s.RecRepo.GetReceptionsByPVZ(ctx, p.ID, startDate, endDate)
		if err != nil {
			continue
		}

		var recWithProducts []*model.ReceptionWithProducts
		for _, r := range receptions {
			products, _ := s.RecRepo.GetProductsByReception(ctx, r.ID)
			recWithProducts = append(recWithProducts, &model.ReceptionWithProducts{
				Reception: r,
				Products:  products,
			})
		}

		result = append(result, &model.PVZWithReceptions{
			PVZ:        p,
			Receptions: recWithProducts,
		})
	}

	offset := (page - 1) * limit
	if offset >= len(result) {
		return []*model.PVZWithReceptions{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], nil
}
