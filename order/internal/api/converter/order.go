package converter

import (
	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/input"
	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
)

func ToCreateOrderInput(req *orderv1.CreateOrderRequest) input.CreateOrderInput {
	var shieldUUID *uuid.UUID
	if req.ShieldUUID.IsSet() {
		v := req.ShieldUUID.Value
		shieldUUID = &v
	}
	var weaponUUID *uuid.UUID
	if req.WeaponUUID.IsSet() {
		v := req.WeaponUUID.Value
		weaponUUID = &v
	}
	return input.CreateOrderInput{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	}
}

func ToOrderDTO(order model.Order) *orderv1.OrderDto {
	var (
		hullUUID   uuid.UUID
		engineUUID uuid.UUID

		shieldUUIDPtr *uuid.UUID
		weaponUUIDPtr *uuid.UUID
	)

	for _, item := range order.Items {
		switch item.PartType {

		case model.PartTypeHull:
			hullUUID = item.PartUUID

		case model.PartTypeEngine:
			engineUUID = item.PartUUID

		case model.PartTypeShield:
			id := item.PartUUID
			shieldUUIDPtr = &id

		case model.PartTypeWeapon:
			id := item.PartUUID
			weaponUUIDPtr = &id
		}
	}

	var shieldUUID orderv1.OptNilUUID
	if shieldUUIDPtr != nil {
		shieldUUID = orderv1.NewOptNilUUID(*shieldUUIDPtr)
	}

	var weaponUUID orderv1.OptNilUUID
	if weaponUUIDPtr != nil {
		weaponUUID = orderv1.NewOptNilUUID(*weaponUUIDPtr)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(
			orderv1.PaymentMethod(*order.PaymentMethod),
		)
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.UUID,
		HullUUID:        hullUUID,
		EngineUUID:      engineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice(),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func ToPaymentMethod(method orderv1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard

	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP

	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard

	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney

	default:
		return model.PaymentMethodUnspecified
	}
}
