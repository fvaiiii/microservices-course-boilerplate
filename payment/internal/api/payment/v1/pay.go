package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/converter"
	errs "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/errors"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

func (s *paymentServer) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	transactionUUID, err := s.paymentService.Pay(ctx, converter.ProtoToModel(req))
	if err != nil {
		if errors.Is(err, errs.ErrInvalidOrderUUID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
