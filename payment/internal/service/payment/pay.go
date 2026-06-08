package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/errors"
	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/service/input"
)

func (s *paymentService) Pay(ctx context.Context, in input.PayOrderInput) (string, error) {
	if in.OrderUUID == "" {
		return "", errs.ErrInvalidOrderUUID
	}
	if _, err := uuid.Parse(in.OrderUUID); err != nil {
		return "", errs.ErrInvalidOrderUUID
	}

	if !in.PaymentMethod.IsValid() {
		return "", errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New().String()

	slog.InfoContext(ctx, "оплата выполнена",
		"order_uuid", in.OrderUUID,
		"transaction_uuid", transactionUUID,
		"payment_method", in.PaymentMethod,
	)

	return transactionUUID, nil
}
