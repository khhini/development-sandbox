package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/khhini/development-sandbox/golang/strategy-pattern/internal/payment"
)

func main() {
	creditCard := &payment.CreditCardStrategy{
		Name:     "Gopher",
		CardNum:  "1234123412341234",
		CVV:      "123",
		ExpMonth: 12,
		ExpYear:  2030,
	}

	paypal := &payment.PayPalStrategy{
		Email:    "gopher@example.com",
		Password: "supersecretpassword",
	}

	crypto := &payment.CryptoStrategy{
		WalletAddress: "0x12342134asf1234",
		CoinType:      "Ethereum",
	}

	processor := payment.NewPaymentProcessor(creditCard)

	var wg sync.WaitGroup
	wg.Add(3)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		processor.SetStrategy(creditCard)
		result, err := processor.ProcessStrategy(100.50)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		processor.SetStrategy(paypal)
		result, err := processor.ProcessStrategy(100.50)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		processor.SetStrategy(crypto)
		result, err := processor.ProcessStrategy(100.50)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}(&wg)

	wg.Wait()
}
