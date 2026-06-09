package idepotency

import (
	"context"
	"time"
)

type IdempotencyRecord struct {
	Key        string            `json:"key"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
	CreatedAt  time.Time         `json:"created_at"`
	ExpiresAt  time.Time         `json:"expires_at"`

	RequestHash string `json:"request_hash"`

	InFlight bool `json:"in_filght"`
}

type IdempotencyStore interface {
	Get(ctx context.Context, key string) (*IdempotencyRecord, error)

	SetInFlight(ctx context.Context, key string, requestHash string, ttl time.Duration) (bool, error)

	Finalize(ctx context.Context, record *IdempotencyRecord) error

	Delete(ctx context.Context, key string) error
}
