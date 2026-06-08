package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/errors"
)

func ErrorInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	switch {
	case errors.Is(err, errs.ErrInvalidOrderUUID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrInvalidPaymentMethod):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
