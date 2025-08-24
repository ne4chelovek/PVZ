package test

import (
	"PVZ/internal/model"
	repoMocks "PVZ/internal/repository/mocks"
	"PVZ/internal/service/reception"
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReceptionService_CreateReception_Success(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := reception.NewReceptionService(receptionMock)

	receptionMock.GetOpenReceptionForPVZMock.Set(func(ctx context.Context, pvzID string) (*model.Reception, error) {
		assert.Equal(t, "1", pvzID)
		return nil, nil
	})

	receptionMock.CreateMock.Set(func(ctx context.Context, reception *model.Reception) (*model.Reception, error) {
		assert.Equal(t, "1", reception.PVZID)
		assert.Equal(t, "in_progress", reception.Status)
		assert.NotZero(t, reception.DateTime)

		return reception, nil
	})

	result, err := svc.CreateReception(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", result.PVZID)
	assert.Equal(t, "in_progress", result.Status)
	assert.NotEmpty(t, result.ID)
	assert.NotZero(t, result.DateTime)
}

func TestReceptionService_CreateReception_AlreadyExists(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)
	svc := reception.NewReceptionService(receptionMock)

	receptionMock.GetOpenReceptionForPVZMock.Set(func(ctx context.Context, pvzID string) (*model.Reception, error) {
		assert.Equal(t, "1", pvzID)
		return &model.Reception{
			ID:     "existing-id",
			PVZID:  "1",
			Status: "in_progress",
		}, nil
	})

	result, err := svc.CreateReception(context.Background(), "1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Nil(t, result)
}
