package input

import "github.com/fvaiiii/microservices-course-boilerplate/payment/internal/model"

type PayOrderInput struct {
	OrderUUID     string
	PaymentMethod model.PaymentMethod
}
