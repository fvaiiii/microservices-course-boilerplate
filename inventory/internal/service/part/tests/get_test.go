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

func TestGet(t *testing.T) {
	t.Parallel()

	var (
		ctx      = context.Background()
		partUUID = uuid.New().String()

		existingPart = model.Part{
			UUID:     partUUID,
			Name:     "Алюминиевый корпус",
			PartType: model.PartTypeHull,
			Price:    500000,
		}
	)

	type expected struct {
		err       error
		onlyError bool
	}

	tests := []struct {
		name      string
		rawUUID   string
		setupMock func(repo *mocks.PartRepository)
		expected  expected
	}{
		{
			name:    "деталь найдена",
			rawUUID: partUUID,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(ctx, partUUID).
					Return(existingPart, nil)
			},
			expected: expected{err: nil},
		},
		{
			name:    "деталь не найдена",
			rawUUID: partUUID,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(ctx, partUUID).
					Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
		{
			name:    "ошибка репозитория",
			rawUUID: partUUID,
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().
					Get(ctx, partUUID).
					Return(model.Part{}, errors.New("ошибка БД"))
			},
			expected: expected{onlyError: true},
		},
		{
			name:      "неверный формат UUID",
			rawUUID:   "not-a-uuid",
			setupMock: func(repo *mocks.PartRepository) {},
			expected:  expected{err: errs.ErrInvalidUUID},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewPartRepository(t)
			tc.setupMock(repo)

			svc := partsvc.New(repo)
			part, err := svc.Get(ctx, tc.rawUUID)

			switch {
			case tc.expected.onlyError:
				require.Error(t, err)
				assert.Empty(t, part.UUID)

			case tc.expected.err != nil:
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, part.UUID)

			default:
				require.NoError(t, err)
				assert.Equal(t, partUUID, part.UUID)
				assert.Equal(t, existingPart.Name, part.Name)
				assert.Equal(t, existingPart.PartType, part.PartType)
				assert.Equal(t, existingPart.Price, part.Price)
			}
		})
	}
}
