package model

import (
	"fmt"
	"time"
)

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"` // Москва, Санкт-Петербург, Казань
}

func (p *PVZ) IsValid() error {
	validCities := map[string]bool{
		"Москва":          true,
		"Санкт-Петербург": true,
		"Казань":          true,
	}
	if !validCities[p.City] {
		return fmt.Errorf("invalid city: %s", p.City)
	}
	return nil
}
