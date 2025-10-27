# Problem

Image you're building an e-commerce system. A user places an order, and you need to:

1. Save the order to your database
2. Send an event to notify other services (inventory, shipping, email, etc.)

```go
func (s *OrderService) CreateOrder(ctx context.Context, order *Order) error {
  if err := s.repo.Save(ctx, order); err != nil {
    return fmt.Errorf("failed to save order: %w", err)
  }

  event := OrderCreatedEvent {
    OrderID: order.ID,
    Ammout: order.TotalAmount,
  }

  // publish event to broker
  if err := s.eventBus.Publish(ctx, "order.created", event); err != nil {
    return fmt.Errorf("failed to publish event: %w", err)
  }

  return nil
} 
```

What could go wrong:

1. The database save succeeds but the message publish fails. You will get an order in your database that no one knows about it.
2. The message publishes successfully, but then something goes wrong and you need to rollback.

# Solution: Transaction Outbox Pattern

Output pattern rely on special "outbox" table in the same database transaction as your business data, instead of publish message directly, the message was written to this special table. All message that will need to be publish are could be reads from this table.

# References

- <https://levelup.gitconnected.com/the-transactional-outbox-pattern-fixing-distributed-message-dd878a407edb>
