package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	inventoryapi "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/api/inventory/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
)

func TestListParts(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		hullUUID    = uuid.New().String()
		engineUUID  = uuid.New().String()
		missingUUID = uuid.New().String()

		allParts = []model.Part{
			{UUID: hullUUID, Name: "Алюминиевый корпус", PartType: model.PartTypeHull, Price: 500000},
			{UUID: engineUUID, Name: "Ионный двигатель C", PartType: model.PartTypeEngine, Price: 300000},
		}

		hullParts = []model.Part{
			{UUID: hullUUID, Name: "Алюминиевый корпус", PartType: model.PartTypeHull, Price: 500000},
		}
	)

	tests := []struct {
		name      string
		req       *inventoryv1.ListPartsRequest
		setupMock func(svc *mocks.PartService)
		wantLen   int
		wantErr   error
	}{
		{
			name: "все детали без фильтра",
			req: &inventoryv1.ListPartsRequest{
				PartType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
			},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					List(ctx, model.PartFilter{PartType: model.PartTypeUnspecified}).
					Return(allParts, nil)
			},
			wantLen: 2,
		},
		{
			name: "фильтр по типу HULL",
			req: &inventoryv1.ListPartsRequest{
				PartType: inventoryv1.PartType_PART_TYPE_HULL,
			},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					List(ctx, model.PartFilter{PartType: model.PartTypeHull}).
					Return(hullParts, nil)
			},
			wantLen: 1,
		},
		{
			name: "фильтр по UUID",
			req: &inventoryv1.ListPartsRequest{
				Uuids: []string{hullUUID, engineUUID},
			},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					List(ctx, model.PartFilter{
						UUIDs:    []string{hullUUID, engineUUID},
						PartType: model.PartTypeUnspecified,
					}).
					Return(allParts, nil)
			},
			wantLen: 2,
		},
		{
			name: "один из UUID не найден",
			req: &inventoryv1.ListPartsRequest{
				Uuids: []string{hullUUID, missingUUID},
			},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					List(ctx, model.PartFilter{
						UUIDs:    []string{hullUUID, missingUUID},
						PartType: model.PartTypeUnspecified,
					}).
					Return(nil, errs.ErrPartNotFound)
			},
			wantErr: errs.ErrPartNotFound,
		},
		{
			name: "неверный формат UUID в фильтре",
			req: &inventoryv1.ListPartsRequest{
				Uuids: []string{hullUUID, "not-a-uuid"},
			},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().
					List(ctx, model.PartFilter{
						UUIDs:    []string{hullUUID, "not-a-uuid"},
						PartType: model.PartTypeUnspecified,
					}).
					Return(nil, errs.ErrInvalidUUID)
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
			res, err := api.ListParts(ctx, tc.req)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			assert.Len(t, res.GetParts(), tc.wantLen)

			if tc.wantLen > 0 {
				assert.Equal(t, hullUUID, res.GetParts()[0].GetUuid())
				assert.Equal(t, inventoryv1.PartType_PART_TYPE_HULL, res.GetParts()[0].GetPartType())
			}
		})
	}
}
