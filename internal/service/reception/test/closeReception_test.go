package test

import (
	"PVZ/internal/model"
	repoMocks "PVZ/internal/repository/mocks"
	"PVZ/internal/service/reception"
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestReceptionService_CloseReception_Success(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := reception.NewReceptionService(receptionMock)

	receptionMock.GetOpenReceptionForPVZMock.Set(func(ctx context.Context, pvzID string) (*model.Reception, error) {
		assert.Equal(t, "1", pvzID)
		return &model.Reception{ID: "1", PVZID: "1", Status: "in_progress", DateTime: time.Now()}, nil
	})

	receptionMock.CloseReceptionMock.Expect(minimock.AnyContext, "1").Return(nil)

	result, err := svc.CloseLastReception(minimock.AnyContext, "1")

	assert.NoError(t, err)
	assert.Equal(t, "close", result.Status)
}

func TestAuthService_CloseReception_Error(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	receptionMock := repoMocks.NewReceptionRepositoryMock(ctrl)

	svc := reception.NewReceptionService(receptionMock)

	receptionMock.GetOpenReceptionForPVZMock.Set(func(ctx context.Context, pvzID string) (*model.Reception, error) {
		return &model.Reception{ID: "1", PVZID: "1", Status: "close"}, nil
	})

	result, err := svc.CloseLastReception(minimock.AnyContext, "1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "нет активной приёмки для закрытия")
}
