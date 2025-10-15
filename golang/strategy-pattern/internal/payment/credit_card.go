package payment

import "fmt"

type CreditCardStrategy struct {
	Name     string
	CardNum  string
	CVV      string
	ExpMonth int
	ExpYear  int
}

func (c *CreditCardStrategy) Process(amount float64) (string, error) {
	return fmt.Sprintf("Paid %.2f using Credit Card ending with %s", amount, c.CardNum[len(c.CardNum)-4:]), nil
}
