package repository

import (
	"PVZ/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type PVZRepository interface {
	Create(ctx context.Context, pvz *model.PVZ) error
	GetAll(ctx context.Context) ([]*model.PVZ, error)
}

type ReceptionRepository interface {
	Create(ctx context.Context, reception *model.Reception) (*model.Reception, error)
	GetOpenReceptionForPVZ(ctx context.Context, pvzID string) (*model.Reception, error)
	CloseReception(ctx context.Context, receptionID string) error
	AddProduct(ctx context.Context, product *model.Product) error
	DeleteLastProductByPVZ(ctx context.Context, pvzID string) error
	GetProductsByReception(ctx context.Context, receptionID string) ([]*model.Product, error)
	GetReceptionsByPVZ(ctx context.Context, pvzID string, startDate, endDate *string) ([]*model.Reception, error)
}
