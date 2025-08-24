package pvz

import (
	"PVZ/internal/model"
	"PVZ/internal/repository"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pvzRepository struct {
	db *pgxpool.Pool
}

func NewPVZRepository(db *pgxpool.Pool) repository.PVZRepository {
	return &pvzRepository{db: db}
}

func (r *pvzRepository) Create(ctx context.Context, pvz *model.PVZ) error {
	query := squirrel.Insert("pvzs").
		Columns("id", "registration_date", "city").
		Values(pvz.ID, pvz.RegistrationDate, pvz.City).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *pvzRepository) GetAll(ctx context.Context) ([]*model.PVZ, error) {
	query := squirrel.Select("id", "registration_date", "city").
		From("pvzs").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pvzs []*model.PVZ
	for rows.Next() {
		var p model.PVZ
		if err := rows.Scan(&p.ID, &p.RegistrationDate, &p.City); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, &p)
	}

	return pvzs, nil
}
