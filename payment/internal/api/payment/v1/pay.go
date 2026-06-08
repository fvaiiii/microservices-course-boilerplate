package v1

import (
	"context"

	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/converter"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

func (s *paymentServer) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	transactionUUID, err := s.paymentService.Pay(ctx, converter.ProtoToModel(req))
	if err != nil {
		return nil, err
	}
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
