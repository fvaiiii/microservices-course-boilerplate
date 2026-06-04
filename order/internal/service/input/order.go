package input

import (
	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
)

type CreateOrderInput struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}
type PayOrderInput struct {
	OrderUUID uuid.UUID
	Method    model.PaymentMethod
}

func (i *CreateOrderInput) PartUUIDs() []uuid.UUID {
	uuids := []uuid.UUID{i.HullUUID, i.EngineUUID}
	if i.ShieldUUID != nil {
		uuids = append(uuids, *i.ShieldUUID)
	}
	if i.WeaponUUID != nil {
		uuids = append(uuids, *i.WeaponUUID)
	}
	return uuids
}
