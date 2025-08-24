package test

import (
	"PVZ/internal/model"
	repoMocks "PVZ/internal/repository/mocks"
	"PVZ/internal/service/pvz"
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPvzService_Create_Success(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	pvzMock := repoMocks.NewPVZRepositoryMock(ctrl)
	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := pvz.NewPVZService(pvzMock, receptionMock)

	pvzMock.CreateMock.Set(func(ctx context.Context, pvz *model.PVZ) error {
		assert.Equal(t, "Москва", pvz.City)
		assert.NotEmpty(t, pvz.ID)
		assert.NotZero(t, pvz.RegistrationDate)

		return nil
	})

	result, err := svc.CreatePVZ(minimock.AnyContext, &model.PVZ{City: "Москва"})

	assert.NoError(t, err)
	assert.Equal(t, "Москва", result.City)
	assert.NotEmpty(t, result.ID)
	assert.NotZero(t, result.RegistrationDate)
}

func TestPvzService_Create_InvalidCity(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	pvzMock := repoMocks.NewPVZRepositoryMock(ctrl)
	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := pvz.NewPVZService(pvzMock, receptionMock)

	result, err := svc.CreatePVZ(minimock.AnyContext, &model.PVZ{City: "Пенза"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Москва, Санкт-Петербург, Казань")
}
