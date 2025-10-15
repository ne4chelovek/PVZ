package test

import (
	"PVZ/internal/service/auth"
	utilsMocks "PVZ/internal/utils/mocks"
	"context"
	"golang.org/x/crypto/bcrypt"
	"testing"

	cacheMocks "PVZ/internal/cache/mocks"
	"PVZ/internal/model"
	repoMocks "PVZ/internal/repository/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login_Success(t *testing.T) {
	t.Parallel()
	// 1. Подготавливаем контроллер
	ctrl := minimock.NewController(t)

	// 2. Мокаем все зависимости
	repoMock := repoMocks.NewUserRepositoryMock(ctrl)
	tokenMock := utilsMocks.NewTokenUtilsMock(ctrl)
	blackListMock := cacheMocks.NewBlackListMock(ctrl)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// 3. Настраиваем ожидания
	expectedUser := &model.User{
		ID:       "123",
		Email:    "test@pvz.com",
		Role:     "employee",
		Password: string(hashedPassword),
	}

	repoMock.GetByEmailMock.Expect(context.Background(), "test@pvz.com").
		Return(expectedUser, nil)

	// Проверяем, что пароль корректный (в реальном коде — bcrypt)
	tokenMock.GenerateTokenMock.Expect("employee").Return("valid-jwt", nil)

	// 4. Создаём сервис
	svc := auth.NewAuthService(blackListMock, tokenMock, repoMock)

	// 5. Вызываем метод
	token, err := svc.Login(context.Background(), "test@pvz.com", "password123")

	// 6. Проверяем
	assert.NoError(t, err)
	assert.Equal(t, "valid-jwt", token)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	repoMock := repoMocks.NewUserRepositoryMock(ctrl)
	tokenMock := utilsMocks.NewTokenUtilsMock(ctrl)
	blackListMock := cacheMocks.NewBlackListMock(ctrl)

	repoMock.GetByEmailMock.Expect(context.Background(), "unknown@pvz.com").
		Return(nil, assert.AnError) // или конкретная ошибка

	svc := auth.NewAuthService(blackListMock, tokenMock, repoMock)

	token, err := svc.Login(context.Background(), "unknown@pvz.com", "pass")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "invalid credentials")
}
