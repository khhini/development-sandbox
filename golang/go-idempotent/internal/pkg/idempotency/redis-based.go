package idepotency

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	keyPrefix   = "idempotency:"
	inFligthTTL = 30 * time.Second
	defaultTTL  = 24 * time.Hour
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (S *RedisStore) redisKey(key string) string {
	return keyPrefix + key
}

func (s *RedisStore) Get(ctx context.Context, key string) (*IdempotencyRecord, error) {
	data, err := s.client.Get(ctx, s.redisKey(key)).Bytes()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var record IdempotencyRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("unmarshal record: %w", err)
	}

	return &record, nil
}

// SetInFlight uses a Lua scirpt for atomic "check-and-set"
// key piece that prevents the race conditions
var setInFlightScript = redis.NewScript(`
	local key = KEYS[1]
	local value = ARGV[1]
	local ttl = tonumber(ARGV[2])

	-- ONly set if key doesn't exist
	local result = redis.call('SET', key, value, 'NX', 'EX', ttl)
	if result then
		return 1 -- Claimed it
	end
	return 0 -- Already exists
	`)

func (s *RedisStore) SetInFlight(ctx context.Context, key string, requestHash string, ttl time.Duration) (bool, error) {
	record := IdempotencyRecord{
		Key:         key,
		InFlight:    true,
		RequestHash: requestHash,
		CreatedAt:   time.Now(),
	}

	data, err := json.Marshal(record)
	if err != nil {
		return false, fmt.Errorf("marshal in-flight record: %w", err)
	}

	result, err := setInFlightScript.Run(ctx, s.client, []string{s.redisKey(key)}, string(data), int(ttl.Seconds())).Int()
	if err != nil {
		return false, fmt.Errorf("set in-flight script: %w", err)
	}

	return result == 1, nil
}

func (s *RedisStore) Finalize(ctx context.Context, record IdempotencyRecord) error {
	record.InFlight = false
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal final record:", err)
	}

	ttl := time.Until(record.ExpiresAt)
	if ttl <= 0 {
		ttl = defaultTTL
	}

	return s.client.Set(ctx, s.redisKey(record.Key), data, ttl).Err()
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.redisKey(key)).Err()
}

func HashRequest(body []byte) string {
	h := sha256.Sum256(body)
	return fmt.Sprintf("%x", h[:8])
}
