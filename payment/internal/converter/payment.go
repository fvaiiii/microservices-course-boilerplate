package converter

import (
	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/payment/internal/service/input"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

func ProtoToModel(req *paymentv1.PayOrderRequest) input.PayOrderInput {
	return input.PayOrderInput{
		OrderUUID:     req.OrderUuid,
		PaymentMethod: PaymentMethodToModel(req.PaymentMethod),
	}
}

func PaymentMethodToModel(paymentMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
