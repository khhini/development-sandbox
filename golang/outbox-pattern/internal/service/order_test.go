package service

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/khhini/development-sandbox/golang/outbox-pattern/internal/domain"
	_ "github.com/lib/pq"
)

func TestOrderService(t *testing.T) {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		t.Errorf("database_url : %v\n", databaseURL)
		t.Errorf("failed to establish database connections: %v", err)
	}

	svc := NewOrderService(db)

	req := domain.CreateOrderRequest{
		CustomerID:  uuid.New(),
		TotalAmount: 100,
		Items:       []string{"Alice", "Bob", "Charlie"},
	}

	t.Run("test create order", func(t *testing.T) {
		_, err := svc.CreateOrder(ctx, req)
		if err != nil {
			t.Errorf("failed create order: %v", err)
		}
	})
}
