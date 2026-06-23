package converter

import (
	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/record"
)

func OrderToRecord(order model.Order) record.Order {
	return record.Order{
		UUID:            order.UUID,
		Status:          string(order.Status),
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethodToRecord(order.PaymentMethod),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderItemsToRecord(orderUUID uuid.UUID, orderItems []model.OrderItem) []record.OrderItem {
	items := make([]record.OrderItem, 0, len(orderItems))

	for _, item := range orderItems {
		items = append(items, record.OrderItem{
			OrderUUID: orderUUID,
			PartUUID:  item.PartUUID,
			PartType:  string(item.PartType),
			Price:     item.Price,
		})
	}
	return items
}

func RecordToModel(order record.Order, items []record.OrderItem) model.Order {
	return model.Order{
		UUID:            order.UUID,
		Items:           OrderItemsToModel(items),
		Status:          model.OrderStatus(order.Status),
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethodToModel(order.PaymentMethod),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderItemsToModel(orderItems []record.OrderItem) []model.OrderItem {
	items := make([]model.OrderItem, 0, len(orderItems))

	for _, item := range orderItems {
		items = append(items, model.OrderItem{
			PartUUID: item.PartUUID,
			PartType: model.PartType(item.PartType),
			Price:    item.Price,
		})
	}
	return items
}

func paymentMethodToRecord(method *model.PaymentMethod) *string {
	if method == nil {
		return nil
	}
	s := string(*method)
	return &s
}

func paymentMethodToModel(method *string) *model.PaymentMethod {
	if method == nil {
		return nil
	}
	m := model.PaymentMethod(*method)
	return &m
}
