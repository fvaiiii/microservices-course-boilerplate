package v1

import (
	"errors"
	"net/http"

	errs "github.com/fvaiiii/microservices-course-boilerplate/order/internal/errors"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func handleCreateOrderError(err error) orderv1.CreateOrderRes {
	switch {
	case errors.Is(err, errs.ErrPartNotFound):
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}

	case errors.Is(err, errs.ErrOutOfStock):
		return &orderv1.CreateOrderConflict{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}

	default:
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
		}
	}
}

func handleGetOrderError(err error) orderv1.GetOrderRes {
	switch {
	case errors.Is(err, errs.ErrOrderNotFound):
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}

	default:
		return &orderv1.GetOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
		}
	}
}

func handlePayOrderError(err error) orderv1.PayOrderRes {
	switch {
	case errors.Is(err, errs.ErrOrderNotFound):
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}

	case errors.Is(err, errs.ErrOrderAlreadyPaid):
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}

	case errors.Is(err, errs.ErrOrderCancelled):
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}

	default:
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
		}
	}
}

func handleCancelOrderError(err error) orderv1.CancelOrderRes {
	switch {
	case errors.Is(err, errs.ErrOrderNotFound):
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}

	case errors.Is(err, errs.ErrOrderAlreadyPaid):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}

	case errors.Is(err, errs.ErrOrderCancelled):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}

	default:
		return &orderv1.CancelOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "internal error",
		}
	}
}
