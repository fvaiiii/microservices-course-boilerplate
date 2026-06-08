package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/client/grpc/payment/v1/converter"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

type client struct {
	api paymentv1.PaymentServiceClient
}

func New(api paymentv1.PaymentServiceClient) *client {
	return &client{
		api: api,
	}
}

func (c *client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	resp, err := c.api.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: converter.ToProtoPaymentMethod(method),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("pay order grpc: %w", err)
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse transaction uuid: %w", err)
	}

	return transactionUUID, nil
}
