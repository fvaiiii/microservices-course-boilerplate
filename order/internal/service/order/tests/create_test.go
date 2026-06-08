package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/input"
	ordersvc "github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		in input.CreateOrderInput
	}

	type expected struct {
		err error
	}

	var (
		ctx = context.Background()

		hullUUID   = uuid.New()
		engineUUID = uuid.New()
		shieldUUID = uuid.New()
		weaponUUID = uuid.New()

		partsInStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", PartType: model.PartTypeHull, Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 5},
		}

		partsOutOfStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", PartType: model.PartTypeHull, Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 0},
		}

		partsWithOptional = []model.Part{
			{UUID: hullUUID, Name: "Hull", PartType: model.PartTypeHull, Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 5},
			{UUID: shieldUUID, Name: "Shield", PartType: model.PartTypeShield, Price: 400000, StockQuantity: 6},
			{UUID: weaponUUID, Name: "Weapon", PartType: model.PartTypeWeapon, Price: 250000, StockQuantity: 7},
		}
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, client *mocks.InventoryClient)
		expected  expected
	}{
		{
			name: "успешное создание заказа",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(partsInStock, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return len(o.Items) == 2 &&
							o.Items[0].PartUUID == hullUUID &&
							o.Items[1].PartUUID == engineUUID &&
							o.TotalPrice() == 800000 &&
							o.Status == model.OrderStatusPendingPayment
					})).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "деталь не найдена",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
		{
			name: "деталь закончилась на складе",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(partsOutOfStock, nil)
			},
			expected: expected{err: errs.ErrOutOfStock},
		},
		{
			name: "успешное создание с опциональными деталями",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
					ShieldUUID: new(shieldUUID),
					WeaponUUID: new(weaponUUID),
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID, shieldUUID, weaponUUID}).
					Return(partsWithOptional, nil)
				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return len(o.Items) == 4 && o.TotalPrice() == 1450000
					})).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "ошибка при сохранении в репозитории",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(partsInStock, nil)
				repo.EXPECT().
					Create(ctx, mock.Anything).
					Return(errors.New("ошибка БД"))
			},
			expected: expected{err: errors.New("ошибка БД")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(orderRepo, inventoryClient)

			svc := ordersvc.New(orderRepo, inventoryClient, paymentClient)
			order, err := svc.Create(ctx, tc.args.in)

			if tc.name == "ошибка при сохранении в репозитории" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "создать заказ")
				assert.Empty(t, order.UUID)
				return
			}

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, order.UUID)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, order.UUID)
			}
		})
	}
}
