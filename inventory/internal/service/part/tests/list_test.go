package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/model"
	partsvc "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/service/part"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/service/part/mocks"
)

func TestList(t *testing.T) {
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

		partsByUUIDs = []model.Part{
			{UUID: hullUUID, Name: "Алюминиевый корпус", PartType: model.PartTypeHull, Price: 500000},
			{UUID: engineUUID, Name: "Ионный двигатель C", PartType: model.PartTypeEngine, Price: 300000},
		}
	)

	type expected struct {
		partsLen  int
		err       error
		onlyError bool
	}

	tests := []struct {
		name      string
		filter    model.PartFilter
		setupMock func(repo *mocks.PartRepository)
		expected  expected
	}{
		{
			name:   "все детали без фильтра",
			filter: model.PartFilter{PartType: model.PartTypeUnspecified},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					List(ctx, model.PartFilter{PartType: model.PartTypeUnspecified}).
					Return(allParts, nil)
			},
			expected: expected{partsLen: 2},
		},
		{
			name:   "фильтр по типу HULL",
			filter: model.PartFilter{PartType: model.PartTypeHull},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					List(ctx, model.PartFilter{PartType: model.PartTypeHull}).
					Return(hullParts, nil)
			},
			expected: expected{partsLen: 1},
		},
		{
			name: "фильтр по UUID",
			filter: model.PartFilter{
				UUIDs: []string{hullUUID, engineUUID},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					List(ctx, model.PartFilter{
						UUIDs: []string{hullUUID, engineUUID},
					}).
					Return(partsByUUIDs, nil)
			},
			expected: expected{partsLen: 2},
		},
		{
			name: "один из UUID не найден",
			filter: model.PartFilter{
				UUIDs: []string{hullUUID, missingUUID},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					List(ctx, model.PartFilter{
						UUIDs: []string{hullUUID, missingUUID},
					}).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
		{
			name: "неверный формат UUID в фильтре",
			filter: model.PartFilter{
				UUIDs: []string{hullUUID, "not-a-uuid"},
			},
			setupMock: func(repo *mocks.PartRepository) {},
			expected:  expected{err: errs.ErrInvalidUUID},
		},
		{
			name:   "ошибка репозитория",
			filter: model.PartFilter{PartType: model.PartTypeEngine},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					List(ctx, model.PartFilter{PartType: model.PartTypeEngine}).
					Return(nil, errors.New("ошибка БД"))
			},
			expected: expected{onlyError: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewPartRepository(t)
			tc.setupMock(repo)

			svc := partsvc.New(repo)
			parts, err := svc.List(ctx, tc.filter)

			switch {
			case tc.expected.onlyError:
				require.Error(t, err)
				assert.Nil(t, parts)

			case tc.expected.err != nil:
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Nil(t, parts)

			default:
				require.NoError(t, err)
				assert.Len(t, parts, tc.expected.partsLen)
			}
		})
	}
}
