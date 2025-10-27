package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/khhini/development-sandbox/golang/outbox-pattern/internal/domain"
)

type OrderService struct {
	db *sqlx.DB
}

func NewOrderService(db *sqlx.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	order := &domain.Order{
		ID:          uuid.New(),
		CustomerID:  req.CustomerID,
		TotalAmount: req.TotalAmount,
		Status:      "PENDING",
		Items:       req.Items,
		CreatedAt:   time.Now(),
	}

	err = s.insertOrder(ctx, tx, order)
	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}

	event := domain.OrderCreatedEvent{
		EventID:     uuid.New(),
		OrderID:     order.ID,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Items:       order.Items,
		Timestamp:   time.Now(),
	}

	err = s.insertOutboxMessage(ctx, tx, "Order", order.ID, "OrderCreated", event)
	if err != nil {
		return nil, fmt.Errorf("insert outbox message: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return order, nil
}

func (s *OrderService) insertOrder(ctx context.Context, tx *sqlx.Tx, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, customer_id, total_amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.ExecContext(ctx, query,
		order.ID,
		order.CustomerID,
		order.TotalAmount,
		order.Status,
		order.CreatedAt,
	)

	return err
}

func (s *OrderService) insertOutboxMessage(
	ctx context.Context,
	tx *sqlx.Tx,
	aggregateType string,
	aggregateID uuid.UUID,
	eventType string,
	payload any,
) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	query := `
		INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = tx.ExecContext(ctx, query,
		uuid.New(),
		aggregateType,
		aggregateID,
		eventType,
		payloadJSON,
		time.Now(),
	)

	return err
}
