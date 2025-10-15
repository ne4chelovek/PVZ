package repository

import (
	"PVZ/internal/model"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryRunner interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type UserRepository interface {
	WithTx(tx pgx.Tx) UserRepository
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type PVZRepository interface {
	WithTx(tx pgx.Tx) PVZRepository
	Create(ctx context.Context, pvz *model.PVZ) error
	GetAll(ctx context.Context) ([]*model.PVZ, error)
}

type ReceptionRepository interface {
	WithTx(tx pgx.Tx) ReceptionRepository
	Create(ctx context.Context, reception *model.Reception) (*model.Reception, error)
	GetOpenReceptionForPVZ(ctx context.Context, pvzID string) (*model.Reception, error)
	CloseReception(ctx context.Context, receptionID string) error
	AddProduct(ctx context.Context, product *model.Product) error
	DeleteLastProductByPVZ(ctx context.Context, pvzID string) error
	GetProductsByReception(ctx context.Context, receptionID string) ([]*model.Product, error)
	GetReceptionsByPVZ(ctx context.Context, pvzID string, startDate, endDate *string) ([]*model.Reception, error)
}

type EventRepository interface {
	WithTx(tx pgx.Tx) EventRepository
	Event(ctx context.Context, event *model.Event) error
	GetUnprocessedEvents(ctx context.Context, tx pgx.Tx, limit int) ([]*model.Event, error)
	MarkAsProcessed(ctx context.Context, tx pgx.Tx, eventID string) error
}
