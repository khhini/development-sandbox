package main

import (
	"errors"
	"fmt"
)

type OrderItem struct{}

type OrderState interface {
	ProcessPayment(o *Order) error
	Ship(o *Order) error
	Cancel(o *Order) error
	Name() string
}

type Order struct {
	id           string
	items        []OrderItem
	totalAmount  float64
	shippingAddr string
	state        OrderState
}

func NewOrder() *Order {
	return &Order{
		state: &PendingState{},
	}
}

func (o *Order) SetState(s OrderState) {
	o.state = s
}

func (o *Order) Status() string {
	return o.state.Name()
}

type PendingState struct{}

func (s *PendingState) Name() string { return "pending" }

func (s *PendingState) ProcessPayment(o *Order) error {
	o.SetState(&PaidState{})
	return nil
}

func (s *PendingState) Ship(o *Order) error {
	return errors.New("payment not processed")
}

func (s *PendingState) Cancel(o *Order) error {
	o.SetState(&CancelledState{})
	return nil
}

type PaidState struct{}

func (s *PaidState) Name() string { return "paid" }

func (s *PaidState) ProcessPayment(o *Order) error {
	return errors.New("order already paid")
}

func (s *PaidState) Ship(o *Order) error {
	o.SetState(&ShippedState{})
	return nil
}

func (s *PaidState) Cancel(o *Order) error {
	o.SetState(&CancelledState{})
	return nil
}

type ShippedState struct{}

func (s *ShippedState) Name() string { return "shipped" }

func (s *ShippedState) ProcessPayment(o *Order) error {
	return errors.New("order already shipped")
}

func (s *ShippedState) Ship(o *Order) error {
	return errors.New("order already shipped")
}

func (s *ShippedState) Cancel(o *Order) error {
	return errors.New("cannot cancel shipped order")
}

type DeliveredState struct{}

func (s *DeliveredState) Name() string { return "delivered" }

func (s *DeliveredState) ProcessPayment(o *Order) error {
	return errors.New("order already delivered")
}

func (s *DeliveredState) Ship(o *Order) error {
	return errors.New("order already delivered")
}

func (s *DeliveredState) Cancel(o *Order) error {
	return errors.New("cannot cancel delivered order")
}

type CancelledState struct{}

func (s *CancelledState) Name() string { return "cancelled" }

func (s *CancelledState) ProcessPayment(o *Order) error {
	return errors.New("cannot process payment for cancelled order")
}

func (s *CancelledState) Ship(o *Order) error {
	return errors.New("cannot ship cancelled order")
}

func (s *CancelledState) Cancel(o *Order) error {
	return errors.New("order already cancelled")
}

func (o *Order) ProcessPayment() error {
	return o.state.ProcessPayment(o)
}

func (o *Order) Ship() error {
	return o.state.Ship(o)
}

func (o *Order) Cancel() error {
	return o.state.Cancel(o)
}

func main() {
	order := NewOrder()

	fmt.Println("Initial status:", order.Status())

	err := order.Ship()
	if err != nil {
		fmt.Println("Ship error:", err)
	}

	err = order.ProcessPayment()
	if err != nil {
		fmt.Println("Payment error:", err)
	}
	fmt.Println("After paid:", order.Status())

	err = order.Ship()
	if err != nil {
		fmt.Println("Ship error:", err)
	}
	fmt.Println("After shipping:", order.Status())

	err = order.Cancel()
	if err != nil {
		fmt.Println("Cancel error:", err)
	}

	order.SetState(&DeliveredState{})
	fmt.Println("Final status:", order.Status())
}
