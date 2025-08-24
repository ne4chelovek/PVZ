package test

import (
	"PVZ/internal/model"
	repoMocks "PVZ/internal/repository/mocks"
	"PVZ/internal/service/pvz"
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPvzService_GetAll(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	pvzMock := repoMocks.NewPVZRepositoryMock(ctrl)
	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := pvz.NewPVZService(pvzMock, receptionMock)

	pvzMock.GetAllMock.Set(func(ctx context.Context) ([]*model.PVZ, error) {
		return []*model.PVZ{{ID: "1", City: "Москва"}}, nil
	})

	receptionMock.GetReceptionsByPVZMock.Set(func(ctx context.Context, pvzID string, startDate, endDate *string) ([]*model.Reception, error) {
		return []*model.Reception{{ID: "1", PVZID: "1", Status: "in_progress"}}, nil
	})

	receptionMock.GetProductsByReceptionMock.Set(func(ctx context.Context, receptionID string) ([]*model.Product, error) {
		return []*model.Product{{ID: "1", Type: "электроника"}}, nil
	})

	result, err := svc.GetAllPVZWithReceptions(minimock.AnyContext, nil, nil, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "1", result[0].PVZ.ID)
	assert.Equal(t, "Москва", result[0].PVZ.City)
	assert.Equal(t, "1", result[0].Receptions[0].Reception.PVZID)
	assert.Equal(t, "in_progress", result[0].Receptions[0].Reception.Status)
	assert.Equal(t, "1", result[0].Receptions[0].Products[0].ID)
	assert.Equal(t, "электроника", result[0].Receptions[0].Products[0].Type)
}

func TestPvzService_GetAll_RepositoryError(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	pvzMock := repoMocks.NewPVZRepositoryMock(ctrl)
	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := pvz.NewPVZService(pvzMock, receptionMock)

	pvzMock.GetAllMock.Set(func(ctx context.Context) ([]*model.PVZ, error) {
		return nil, fmt.Errorf("failed to query PVZ from database")
	})

	result, err := svc.GetAllPVZWithReceptions(context.Background(), nil, nil, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get PVZ")
}
