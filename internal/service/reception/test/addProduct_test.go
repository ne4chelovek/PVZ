package test

import (
	"PVZ/internal/model"
	repoMocks "PVZ/internal/repository/mocks"
	"PVZ/internal/service/pvz"
	"PVZ/internal/service/reception"
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestReceptionService_AddProduct_Success(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := reception.NewReceptionService(receptionMock)

	receptionMock.GetOpenReceptionForPVZMock.Set(func(ctx context.Context, pvzID string) (*model.Reception, error) {
		assert.Equal(t, "1", pvzID)
		return &model.Reception{ID: "1", PVZID: "1", Status: "in_progress", DateTime: time.Now()}, nil
	})

	receptionMock.AddProductMock.Set(func(ctx context.Context, product *model.Product) error {
		assert.Equal(t, "1", product.ReceptionID)
		assert.Equal(t, "электроника", product.Type)
		assert.NotEmpty(t, product.ID)
		assert.NotZero(t, product.DateTime)
		return nil
	})

	result, err := svc.AddProduct(minimock.AnyContext, "1", "электроника")

	assert.NoError(t, err)
	assert.Equal(t, "электроника", result.Type)
	assert.NotEmpty(t, result.ID)
	assert.NotZero(t, result.DateTime)
}

func TestAuthService_Create_InvalidCity(t *testing.T) {
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
