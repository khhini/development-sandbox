package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID `json:"id"`
	CustomerID  uuid.UUID `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	Items       []string  `json:"items"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderCreatedEvent struct {
	EventID     uuid.UUID `json:"event_id"`
	OrderID     uuid.UUID `json:"order_id"`
	CustomerID  uuid.UUID `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	Items       []string  `json:"items"`
	Timestamp   time.Time `json:"timestamp"`
}

type CreateOrderRequest struct {
	CustomerID  uuid.UUID `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	Items       []string  `json:"items"`
}
