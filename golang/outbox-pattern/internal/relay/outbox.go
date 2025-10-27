package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/khhini/development-sandbox/golang/outbox-pattern/internal/domain"
)

type MessagePublisher interface {
	Publish(ctx context.Context, eventType string, payload []byte) error
}

type OutboxRelay struct {
	db           *sqlx.DB
	publisher    MessagePublisher
	pollInterval time.Duration
	batchSize    int
}

func NewOutboxRelay(db *sqlx.DB, publisher MessagePublisher) *OutboxRelay {
	return &OutboxRelay{
		db:           db,
		publisher:    publisher,
		pollInterval: 1 * time.Second,
		batchSize:    100,
	}
}

func (r *OutboxRelay) Start(ctx context.Context) error {
	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := r.processMessage(ctx); err != nil {
				log.Printf("Error processing outbox: %v", err)
			}
		}
	}
}

func (r *OutboxRelay) processMessage(ctx context.Context) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	messages, err := r.fetchUnprocessedMessages(ctx, tx)
	if err != nil {
		return fmt.Errorf("fetch message: %w", err)
	}

	if len(messages) == 0 {
		return nil
	}

	processed := 0
	for _, msg := range messages {
		if err := r.publishAndMarkProcessed(ctx, tx, msg); err != nil {
			log.Printf("Failed to process message %s: %v", msg.ID, err)
			continue
		}
		processed++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *OutboxRelay) fetchUnprocessedMessages(ctx context.Context, tx *sqlx.Tx) ([]domain.OutboxMessage, error) {
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, created_at
		FROM outbox
		WHERE processed_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := tx.QueryContext(ctx, query, r.batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.OutboxMessage
	for rows.Next() {
		var msg domain.OutboxMessage
		err := rows.Scan(
			&msg.ID,
			&msg.AggregateType,
			&msg.AggregateID,
			&msg.EventType,
			&msg.Payload,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *OutboxRelay) publishAndMarkProcessed(ctx context.Context, tx *sqlx.Tx, msg domain.OutboxMessage) error {
	if err := r.publisher.Publish(ctx, msg.EventType, []byte(msg.Payload)); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	query := `UPDATE outbox SET processed_at = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, time.Now(), msg.ID)
	if err != nil {
		return fmt.Errorf("update outbox: %w", err)
	}

	log.Printf("Publish message: %s (type: %s, aggregate: %s)",
		msg.ID, msg.EventType, msg.AggregateID)

	return nil
}
