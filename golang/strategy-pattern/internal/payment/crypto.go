package payment

import "fmt"

type CryptoStrategy struct {
	WalletAddress string
	CoinType      string
}

func (c *CryptoStrategy) Process(amount float64) (string, error) {
	return fmt.Sprintf("Paid %.2f worth of %s to wallet: %s", amount, c.CoinType, c.WalletAddress), nil
}
