package test

import (
	"PVZ/internal/model"
	repoMocks "PVZ/internal/repository/mocks"
	"PVZ/internal/service/reception"
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReceptionService_DeleteLastProduct_Success(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := reception.NewReceptionService(receptionMock)

	receptionMock.GetOpenReceptionForPVZMock.Set(func(ctx context.Context, pvzID string) (*model.Reception, error) {
		assert.Equal(t, "1", pvzID)
		return &model.Reception{ID: "1", DateTime: time.Now(), PVZID: pvzID, Status: "in_progress"}, nil
	})

	receptionMock.GetProductsByReceptionMock.Set(func(ctx context.Context, receptionID string) (ppa1 []*model.Product, err error) {
		assert.Equal(t, "1", receptionID)
		return []*model.Product{{ID: "1", DateTime: time.Now(), Type: "электроника", ReceptionID: receptionID}}, nil
	})

	receptionMock.DeleteLastProductByPVZMock.Set(func(ctx context.Context, pvzID string) error {
		assert.Equal(t, "1", pvzID)
		return nil
	})

	err := svc.DeleteLastProduct(minimock.AnyContext, "1")

	assert.NoError(t, err)
}

func TestReceptionService_DeleteLastProduct_GetReceptionRepoError(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)
	svc := reception.NewReceptionService(receptionMock)

	receptionMock.GetOpenReceptionForPVZMock.Set(func(ctx context.Context, pvzID string) (*model.Reception, error) {
		assert.Equal(t, "1", pvzID)
		return nil, fmt.Errorf("database error")
	})

	err := svc.DeleteLastProduct(context.Background(), "1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "нет активной приёмки")
	assert.Contains(t, err.Error(), "database error")
}
