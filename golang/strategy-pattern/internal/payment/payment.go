package payment

type PaymentStrategy interface {
	Process(amount float64) (string, error)
}

type PaymentProcessor struct {
	strategy PaymentStrategy
}

func NewPaymentProcessor(strategy PaymentStrategy) *PaymentProcessor {
	return &PaymentProcessor{
		strategy,
	}
}

func (p *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
	p.strategy = strategy
}

func (p *PaymentProcessor) ProcessStrategy(amount float64) (string, error) {
	return p.strategy.Process(amount)
}
