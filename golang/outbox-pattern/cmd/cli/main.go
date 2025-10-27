package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/khhini/development-sandbox/golang/outbox-pattern/internal/domain"
	"github.com/khhini/development-sandbox/golang/outbox-pattern/internal/publisher"
	"github.com/khhini/development-sandbox/golang/outbox-pattern/internal/relay"
	"github.com/khhini/development-sandbox/golang/outbox-pattern/internal/service"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	publisher, _ := publisher.NewMockPublisher()

	relay := relay.NewOutboxRelay(db, publisher)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := relay.Start(ctx); err != nil {
			log.Printf("relay error: %v", err)
		}
	}()

	orderService := service.NewOrderService(db)
	order, err := orderService.CreateOrder(ctx, domain.CreateOrderRequest{
		CustomerID:  uuid.New(),
		TotalAmount: 149.99,
		Items:       []string{"laptop", "mouse"},
	})

	if err != nil {
		log.Printf("failed to create order: %v", err)
	} else {
		log.Printf("order created successfully: %s", order.ID)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	time.Sleep(2 * time.Second)
}
