package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

func mapCreateOrderInventoryError(err error) orderv1.CreateOrderRes {
	st, ok := status.FromError(err)
	if !ok {
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка inventory сервиса",
		}
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: st.Message(),
		}
	case codes.NotFound:
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: st.Message(),
		}
	default:
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка inventory сервиса",
		}
	}
}

func mapHTTPPaymentMethodToGRPC(method orderv1.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch method {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, errors.New("payment_method должен быть указан")
	}
}

func mapPayOrderPaymentError(err error) orderv1.PayOrderRes {
	st, ok := status.FromError(err)
	if !ok {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка payment сервиса",
		}
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: st.Message(),
		}
	case codes.NotFound:
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: st.Message(),
		}
	default:
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка payment сервиса",
		}
	}
}

func collectCreateOrderPartUUIDs(
	req *orderv1.CreateOrderRequest,
) ([]string, orderv1.OptNilUUID, orderv1.OptNilUUID, orderv1.CreateOrderRes) {
	partUUIDs := make([]string, 0, 4)
	partUUIDs = append(partUUIDs, req.GetHullUUID().String())
	partUUIDs = append(partUUIDs, req.GetEngineUUID().String())

	shieldOpt := req.GetShieldUUID()
	if shieldOpt.Set && shieldOpt.Null {
		return nil, orderv1.OptNilUUID{}, orderv1.OptNilUUID{}, &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "shield_uuid не может быть null",
		}
	}
	if shieldOpt.Set && !shieldOpt.Null {
		partUUIDs = append(partUUIDs, shieldOpt.Value.String())
	}

	weaponOpt := req.GetWeaponUUID()
	if weaponOpt.Set && weaponOpt.Null {
		return nil, orderv1.OptNilUUID{}, orderv1.OptNilUUID{}, &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "weapon_uuid не может быть null",
		}
	}
	if weaponOpt.Set && !weaponOpt.Null {
		partUUIDs = append(partUUIDs, weaponOpt.Value.String())
	}

	return partUUIDs, shieldOpt, weaponOpt, nil
}

func (h *OrderHandler) loadCreateOrderParts(
	ctx context.Context,
	partUUIDs []string,
) (*inventoryv1.ListPartsResponse, orderv1.CreateOrderRes) {
	grpcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.inventoryClient.ListParts(grpcCtx, &inventoryv1.ListPartsRequest{
		Uuids: partUUIDs,
	})
	if err != nil {
		return nil, mapCreateOrderInventoryError(err)
	}

	if len(resp.GetParts()) != len(partUUIDs) {
		return nil, &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "не удалось получить все детали",
		}
	}

	return resp, nil
}

func validateCreateOrderParts(
	req *orderv1.CreateOrderRequest,
	resp *inventoryv1.ListPartsResponse,
	shieldOpt orderv1.OptNilUUID,
	weaponOpt orderv1.OptNilUUID,
) (int64, *uuid.UUID, *uuid.UUID, orderv1.CreateOrderRes) {
	var (
		totalPrice int64
		shieldUUID *uuid.UUID
		weaponUUID *uuid.UUID
	)

	for _, part := range resp.GetParts() {
		partPrice, partShieldUUID, partWeaponUUID, badRes := validateSingleCreateOrderPart(req, part, shieldOpt, weaponOpt)
		if badRes != nil {
			return 0, nil, nil, badRes
		}

		if partShieldUUID != nil {
			shieldUUID = partShieldUUID
		}

		if partWeaponUUID != nil {
			weaponUUID = partWeaponUUID
		}

		totalPrice += partPrice
	}

	return totalPrice, shieldUUID, weaponUUID, nil
}

func validateSingleCreateOrderPart(
	req *orderv1.CreateOrderRequest,
	part *inventoryv1.Part,
	shieldOpt orderv1.OptNilUUID,
	weaponOpt orderv1.OptNilUUID,
) (int64, *uuid.UUID, *uuid.UUID, orderv1.CreateOrderRes) {
	hasShield := shieldOpt.Set && !shieldOpt.Null
	hasWeapon := weaponOpt.Set && !weaponOpt.Null

	if part.GetStockQuantity() <= 0 {
		return 0, nil, nil, &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("деталь %s отсутствует на складе", part.GetName()),
		}
	}

	switch part.GetUuid() {
	case req.GetHullUUID().String():
		if part.GetPartType() != inventoryv1.PartType_PART_TYPE_HULL {
			return 0, nil, nil, &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "hull_uuid должен указывать на корпус",
			}
		}
		return part.GetPrice(), nil, nil, nil

	case req.GetEngineUUID().String():
		if part.GetPartType() != inventoryv1.PartType_PART_TYPE_ENGINE {
			return 0, nil, nil, &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "engine_uuid должен указывать на двигатель",
			}
		}
		return part.GetPrice(), nil, nil, nil

	case shieldOpt.Value.String():
		if hasShield {
			if part.GetPartType() != inventoryv1.PartType_PART_TYPE_SHIELD {
				return 0, nil, nil, &orderv1.CreateOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: "shield_uuid должен указывать на щит",
				}
			}
			id := shieldOpt.Value
			return part.GetPrice(), &id, nil, nil
		}

	case weaponOpt.Value.String():
		if hasWeapon {
			if part.GetPartType() != inventoryv1.PartType_PART_TYPE_WEAPON {
				return 0, nil, nil, &orderv1.CreateOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: "weapon_uuid должен указывать на оружие",
				}
			}
			id := weaponOpt.Value
			return part.GetPrice(), nil, &id, nil
		}
	}

	return part.GetPrice(), nil, nil, nil
}
