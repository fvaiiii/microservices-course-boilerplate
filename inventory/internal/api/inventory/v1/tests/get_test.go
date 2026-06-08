package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	inventoryapi "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/api/inventory/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

func TestGetPart(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		partUUID  = uuid.New().String()
		createdAt = time.Now()

		existingPart = model.Part{
			UUID:          partUUID,
			Name:          "Алюминиевый корпус",
			Description:   "Лёгкий корпус",
			Price:         500000,
			PartType:      model.PartTypeHull,
			StockQuantity: 10,
			CreatedAt:     createdAt,
		}
	)

	tests := []struct {
		name      string
		req       *inventoryv1.GetPartRequest
		setupMock func(svc *mocks.PartService)
		wantErr   error
	}{
		{
			name: "успешное получение детали",
			req:  &inventoryv1.GetPartRequest{Uuid: partUUID},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					Get(ctx, partUUID).
					Return(existingPart, nil)
			},
		},
		{
			name: "деталь не найдена",
			req:  &inventoryv1.GetPartRequest{Uuid: partUUID},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					Get(ctx, partUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
		{
			name: "неверный формат UUID",
			req:  &inventoryv1.GetPartRequest{Uuid: "not-a-uuid"},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					Get(ctx, "not-a-uuid").
					Return(model.Part{}, errs.ErrInvalidUUID)
			},
			wantErr: errs.ErrInvalidUUID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewPartService(t)
			tc.setupMock(svc)

			api := inventoryapi.New(svc)
			res, err := api.GetPart(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)

			part := res.GetPart()
			require.NotNil(t, part)
			assert.Equal(t, partUUID, part.GetUuid())
			assert.Equal(t, existingPart.Name, part.GetName())
			assert.Equal(t, existingPart.Description, part.GetDescription())
			assert.Equal(t, existingPart.Price, part.GetPrice())
			assert.Equal(t, inventoryv1.PartType_PART_TYPE_HULL, part.GetPartType())
			assert.Equal(t, existingPart.StockQuantity, part.GetStockQuantity())
			assert.NotNil(t, part.GetCreatedAt())
		})
	}
}
