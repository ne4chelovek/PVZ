package handler

import (
	"PVZ/internal/service"
)

type PVZHandler struct {
	pvzService       service.PVZService
	receptionService service.ReceptionService
	authService      service.AuthService
}

func NewPVZHandler(pvzService service.PVZService, receptionService service.ReceptionService, authService service.AuthService) *PVZHandler {
	return &PVZHandler{
		pvzService:       pvzService,
		receptionService: receptionService,
		authService:      authService,
	}
}
