package payment

import "fmt"

type PayPalStrategy struct {
	Email    string
	Password string
}

func (p *PayPalStrategy) Process(ammount float64) (string, error) {
	return fmt.Sprintf("Paid %.2f using PayPal account: %s", ammount, p.Email), nil
}
