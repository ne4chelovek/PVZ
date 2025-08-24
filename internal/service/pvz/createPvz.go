package pvz

import (
	"PVZ/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
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

	if err := s.Repo.Create(ctx, pvz); err != nil {
		return nil, fmt.Errorf("failed to create PVZ in repository: %w", err)
	}

	return pvz, nil
}
