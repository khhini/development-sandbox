package publisher

import (
	"context"
	"fmt"
)

type MockPublisher struct{}

func NewMockPublisher() (*MockPublisher, error) {
	return &MockPublisher{}, nil
}

func (p *MockPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	msg := string(payload)

	fmt.Printf("eventType: %s", eventType)
	fmt.Printf("msg: %s", msg)

	return nil
}

func (p *MockPublisher) Close() error {
	return nil
}
