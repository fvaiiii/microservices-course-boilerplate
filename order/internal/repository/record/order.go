package record

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID            uuid.UUID
	Items           []OrderItem
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string
	CreatedAt       time.Time
}

type OrderItem struct {
	PartUUID uuid.UUID
	PartType string
	Price    int64
}
