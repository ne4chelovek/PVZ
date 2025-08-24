package service

import (
	"PVZ/internal/model"
	"context"
)

type AuthService interface {
	LoginDummy(ctx context.Context, role string) (string, error)
	Register(ctx context.Context, email, password, role string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type PVZService interface {
	CreatePVZ(ctx context.Context, pvz *model.PVZ) (*model.PVZ, error)
	GetAllPVZWithReceptions(ctx context.Context, startDate, endDate *string, page, limit int) ([]*model.PVZWithReceptions, error)
}

type ReceptionService interface {
	CreateReception(ctx context.Context, pvzID string) (*model.Reception, error)
	AddProduct(ctx context.Context, pvzID string, productType string) (*model.Product, error)
	CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error)
	DeleteLastProduct(ctx context.Context, pvzID string) error
}
