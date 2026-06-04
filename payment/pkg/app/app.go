package app

import (
	"google.golang.org/grpc"

	paymentapi "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/api/payment/v1"
	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/interceptor"
	paymentsvc "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/service/payment"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

func RegisterServices(grpcServer *grpc.Server) {
	svc := paymentsvc.New()
	api := paymentapi.New(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
}
