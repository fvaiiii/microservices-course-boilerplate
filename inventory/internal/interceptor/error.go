package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/errors"
)

func ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	switch {
	case errors.Is(err, errs.ErrInvalidUUID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrPartNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
