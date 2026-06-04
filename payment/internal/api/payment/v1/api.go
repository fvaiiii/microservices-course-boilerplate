package v1

import (
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

type paymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

func New(paymentService PaymentService) *paymentServer {
	return &paymentServer{
		paymentService: paymentService,
	}
}
