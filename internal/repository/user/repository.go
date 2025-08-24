package repository

import (
	"PVZ/internal/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
	sq squirrel.StatementBuilderType
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := r.sq.Insert("users").
		Columns("id", "email", "password", "role").
		Values(user.ID, user.Email, user.Password, user.Role).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build Create user query: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute Create user query: %w", err)
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := r.sq.Select("id", "email", "password", "role").
		From("users").
		Where(squirrel.Eq{"email": email}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build GetByEmail query: %w", err)
	}

	var user model.User
	err = r.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}
