package test

import (
	"PVZ/internal/model"
	"PVZ/internal/service/auth"
	utilsMocks "PVZ/internal/utils/mocks"
	"context"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"

	repoMocks "PVZ/internal/repository/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Register_Success(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	tokenMock := utilsMocks.NewTokenUtilsMock(ctrl)
	repoMock := repoMocks.NewUserRepositoryMock(ctrl)

	svc := auth.NewAuthService(nil, tokenMock, repoMock)

	repoMock.CreateMock.Set(func(ctx context.Context, user *model.User) error {
		assert.Equal(t, "test@pvz.com", user.Email)
		assert.Equal(t, "employee", user.Role)
		assert.NotEmpty(t, user.ID)

		assert.NotEqual(t, "password", user.Password)
		assert.True(t, len(user.Password) > 0)

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
		assert.NoError(t, err, "Hashed password is invalid")

		return nil
	})

	user, err := svc.Register(context.Background(), "test@pvz.com", "password", "employee")

	assert.NoError(t, err)
	assert.Equal(t, "test@pvz.com", user.Email)
	assert.Equal(t, "employee", user.Role)
	assert.NotEmpty(t, user.ID)
}

func TestAuthService_Register_FullInvalidInput(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	tokenMock := utilsMocks.NewTokenUtilsMock(ctrl)
	repoMock := repoMocks.NewUserRepositoryMock(ctrl)

	svc := auth.NewAuthService(nil, tokenMock, repoMock)

	user, err := svc.Register(minimock.AnyContext, "", "", "")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, strings.Contains(err.Error(), "email") ||
		strings.Contains(err.Error(), "password") ||
		strings.Contains(err.Error(), "role"))

}

func TestAuthService_Register_InvalidRole(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	tokenMock := utilsMocks.NewTokenUtilsMock(ctrl)
	repoMock := repoMocks.NewUserRepositoryMock(ctrl)

	svc := auth.NewAuthService(nil, tokenMock, repoMock)

	user, err := svc.Register(minimock.AnyContext, "test@pvz.com", "password", "nil")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid role")
}
