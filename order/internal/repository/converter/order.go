package converter

import (
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/model"
	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/record"
)

func OrderToRecord(order model.Order) record.Order {
	return record.Order{
		UUID:            order.UUID,
		Items:           OrderItemsToRecord(order.Items),
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   (*string)(order.PaymentMethod),
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderItemsToRecord(orderItems []model.OrderItem) []record.OrderItem {
	items := make([]record.OrderItem, 0, len(orderItems))

	for _, item := range orderItems {
		items = append(items, record.OrderItem{
			PartUUID: item.PartUUID,
			PartType: string(item.PartType),
			Price:    item.Price,
		})
	}
	return items
}

func RecordToModel(order record.Order) model.Order {
	return model.Order{
		UUID:            order.UUID,
		Items:           OrderItemsToModel(order.Items),
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   (*model.PaymentMethod)(order.PaymentMethod),
		Status:          model.OrderStatus(order.Status),
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
