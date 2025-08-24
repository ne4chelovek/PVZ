package repository

import (
	"PVZ/internal/model"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Masterminds/squirrel"
)

type receptionRepository struct {
	db *pgxpool.Pool
}

func NewReceiptRepository(db *pgxpool.Pool) *receptionRepository {
	return &receptionRepository{db: db}
}

func (r *receptionRepository) Create(ctx context.Context, reception *model.Reception) (*model.Reception, error) {
	query := squirrel.Insert("receptions").
		Columns("id", "date_time", "pvz_id", "status").
		Values(reception.ID, reception.DateTime, reception.PVZID, reception.Status).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build Create query: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Create query: %w", err)
	}

	return reception, nil
}

func (r *receptionRepository) GetOpenReceptionForPVZ(ctx context.Context, pvzID string) (*model.Reception, error) {
	query := squirrel.Select("id", "date_time", "pvz_id", "status").
		From("receptions").
		Where(squirrel.Eq{
			"pvz_id": pvzID,
			"status": "in_progress",
		}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build GetOpenReceptionForPVZ query: %w", err)
	}

	var rec model.Reception
	err = r.db.QueryRow(ctx, sql, args...).Scan(&rec.ID, &rec.DateTime, &rec.PVZID, &rec.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &rec, nil
}

func (r *receptionRepository) CloseReception(ctx context.Context, receptionID string) error {
	query := squirrel.Update("receptions").
		Set("status", "close").
		Where(squirrel.Eq{"id": receptionID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build CloseReception query: %w", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute CloseReception: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reception with id %s not found or already closed", receptionID)
	}

	return nil
}

func (r *receptionRepository) AddProduct(ctx context.Context, product *model.Product) error {
	query := squirrel.Insert("products").
		Columns("id", "date_time", "type", "reception_id").
		Values(product.ID, product.DateTime, product.Type, product.ReceptionID).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build AddProduct query: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute AddProduct: %w", err)
	}

	return nil
}

func (r *receptionRepository) DeleteLastProductByPVZ(ctx context.Context, pvzID string) error {
	// Шаг 1: Найти ID последнего товара в активной приёмке
	query := squirrel.Select("p.id").
		From("products p").
		Join("receptions r ON p.reception_id = r.id").
		Where(squirrel.Eq{
			"r.pvz_id": pvzID,
			"r.status": "in_progress",
		}).
		OrderBy("p.date_time DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build DeleteLastProductByPVZ select query: %w", err)
	}

	var productID string
	err = r.db.QueryRow(ctx, sql, args...).Scan(&productID)
	if err != nil {
		return fmt.Errorf("failed to get last product: %w", err)
	}

	// Шаг 2: Удалить товар
	deleteQuery := squirrel.Delete("products").
		Where(squirrel.Eq{"id": productID})

	sql, args, err = deleteQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build DeleteLastProductByPVZ delete query: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *receptionRepository) GetProductsByReception(ctx context.Context, receptionID string) ([]*model.Product, error) {
	query := squirrel.Select("id", "date_time", "type", "reception_id").
		From("products").
		Where(squirrel.Eq{"reception_id": receptionID}).
		OrderBy("date_time ASC").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build GetProductsByReception query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.DateTime, &p.Type, &p.ReceptionID)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		products = append(products, &p)
	}

	return products, nil
}

func (r *receptionRepository) GetReceptionsByPVZ(ctx context.Context, pvzID string, startDate, endDate *string) ([]*model.Reception, error) {
	query := squirrel.Select("id", "date_time", "pvz_id", "status").
		From("receptions").
		Where(squirrel.Eq{"pvz_id": pvzID}).
		PlaceholderFormat(squirrel.Dollar)

	// Фильтрация по дате
	if startDate != nil {
		query = query.Where(squirrel.Gt{"date_time": *startDate})
	}
	if endDate != nil {
		query = query.Where(squirrel.Lt{"date_time": *endDate})
	}

	//Мб так если с gt,lt не сработает
	//if startDate != nil {
	//	query = query.Where(squirrel.Expr("date_time >= ?", *startDate))
	//}
	//if endDate != nil {
	//	query = query.Where(squirrel.Expr("date_time <= ?", *endDate))
	//}

	query = query.OrderBy("date_time DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build GetReceptionsByPVZ query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var receptions []*model.Reception
	for rows.Next() {
		var r model.Reception
		err := rows.Scan(&r.ID, &r.DateTime, &r.PVZID, &r.Status)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		receptions = append(receptions, &r)
	}

	return receptions, nil
}
