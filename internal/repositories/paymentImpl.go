package repositories

import (
	"fmt"
	"pattern/internal/models"
)

// https://github.com/abdymazhit/kaspi-merchant-api
type KaspiPayment struct{}

// ProcessPayment(payment *models.Payment) error
func (p *KaspiPayment) ProcessPayment(payment *models.Payment) error {
	fmt.Println("kaspi")
	return nil
}

// https://github.com/plutov/paypal
type PayPalPayment struct{}

func NewPayPalPayment() *PayPalPayment {
	return &PayPalPayment{}
}

func (p *PayPalPayment) ProcessPayment(payment *models.Payment) error {
	fmt.Println("paypal")
	return nil
}
