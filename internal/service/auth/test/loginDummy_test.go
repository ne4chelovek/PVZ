package test

import (
	"PVZ/internal/service/auth"
	"testing"

	utilsMocks "PVZ/internal/utils/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_LoginDummy_Success(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	tokenMock := utilsMocks.NewTokenUtilsMock(ctrl)

	tokenMock.GenerateTokenMock.Expect("employee").Return("valid-jwt", nil)

	svc := auth.NewAuthService(nil, tokenMock, nil)

	token, err := svc.LoginDummy(minimock.AnyContext, "employee")

	assert.NoError(t, err)
	assert.Equal(t, "valid-jwt", token)
}

func TestAuthService_LoginDummy_InvalidRole(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	tokenMock := utilsMocks.NewTokenUtilsMock(ctrl)
	svc := auth.NewAuthService(nil, tokenMock, nil)

	token, err := svc.LoginDummy(minimock.AnyContext, "hacker")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "invalid role")
}
