package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/record"
	serviceorder "github.com/fvaiiii/microservices-course-boilerplate/order/internal/service/order"
)

type repository struct {
	data map[uuid.UUID]record.Order
	mu   sync.RWMutex
}

func New() serviceorder.OrderRepository {
	return &repository{
		data: make(map[uuid.UUID]record.Order),
	}
}
