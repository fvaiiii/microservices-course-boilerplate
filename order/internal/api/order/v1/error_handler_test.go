package v1

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func TestHandleCreateOrderError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "деталь не найдена",
			err:        errs.ErrPartNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "обёрнутая ошибка деталь не найдена",
			err:        fmt.Errorf("получить детали: %w", errs.ErrPartNotFound),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "деталь отсутствует на складе",
			err:        errs.ErrOutOfStock,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "внутренняя ошибка",
			err:        errors.New("неожиданная ошибка"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := handleCreateOrderError(tc.err)
			assertHTTPCode(t, res, tc.wantStatus)
		})
	}
}

func TestHandleGetOrderError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "заказ не найден",
			err:        errs.ErrOrderNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "внутренняя ошибка",
			err:        errors.New("неожиданная ошибка"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := handleGetOrderError(tc.err)
			assertHTTPCode(t, res, tc.wantStatus)
		})
	}
}

func TestHandlePayOrderError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "заказ не найден",
			err:        errs.ErrOrderNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "заказ уже оплачен",
			err:        errs.ErrOrderAlreadyPaid,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "заказ отменён",
			err:        errs.ErrOrderCancelled,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "внутренняя ошибка",
			err:        errors.New("неожиданная ошибка"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := handlePayOrderError(tc.err)
			assertHTTPCode(t, res, tc.wantStatus)
		})
	}
}

func TestHandleCancelOrderError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "заказ не найден",
			err:        errs.ErrOrderNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "заказ уже оплачен",
			err:        errs.ErrOrderAlreadyPaid,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "заказ уже отменён",
			err:        errs.ErrOrderCancelled,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "внутренняя ошибка",
			err:        errors.New("неожиданная ошибка"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := handleCancelOrderError(tc.err)
			assertHTTPCode(t, res, tc.wantStatus)
		})
	}
}

func assertHTTPCode(t *testing.T, res any, wantStatus int) {
	t.Helper()

	switch r := res.(type) {
	case *orderv1.CreateOrderNotFound:
		assert.Equal(t, wantStatus, r.Code)
	case *orderv1.CreateOrderConflict:
		assert.Equal(t, wantStatus, r.Code)
	case *orderv1.CreateOrderInternalServerError:
		assert.Equal(t, wantStatus, r.Code)

	case *orderv1.GetOrderNotFound:
		assert.Equal(t, wantStatus, r.Code)
	case *orderv1.GetOrderInternalServerError:
		assert.Equal(t, wantStatus, r.Code)

	case *orderv1.PayOrderNotFound:
		assert.Equal(t, wantStatus, r.Code)
	case *orderv1.PayOrderConflict:
		assert.Equal(t, wantStatus, r.Code)
	case *orderv1.PayOrderInternalServerError:
		assert.Equal(t, wantStatus, r.Code)

	case *orderv1.CancelOrderNotFound:
		assert.Equal(t, wantStatus, r.Code)
	case *orderv1.CancelOrderConflict:
		assert.Equal(t, wantStatus, r.Code)
	case *orderv1.CancelOrderInternalServerError:
		assert.Equal(t, wantStatus, r.Code)

	default:
		require.Fail(t, "неизвестный тип ответа", "%T", res)
	}
}
