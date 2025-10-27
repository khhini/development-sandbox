package domain

import (
	"time"

	"github.com/google/uuid"
)

type OutboxMessage struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	AggregateType string     `db:"aggregate_type" json:"aggregate_type"`
	AggregateID   uuid.UUID  `db:"aggregate_id" json:"aggregate_id"`
	EventType     string     `db:"event_type" json:"event_type"`
	Payload       string     `db:"payload" json:"payload"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	ProcessedAt   *time.Time `db:"processed_at" json:"processed_at,omitempty"`
}
