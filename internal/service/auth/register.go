package auth

import (
	"PVZ/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(ctx context.Context, email, password, role string) (*model.User, error) {
	if !validRoles[role] {
		return nil, fmt.Errorf("invalid role")
	}

	// Хешируем пароль
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(hashed),
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
