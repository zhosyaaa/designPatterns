package repositories

import (
	"fmt"
	"pattern/internal/models"
)

type KaspiPayment struct{}

// ProcessPayment(payment *models.Payment) error
func (p *KaspiPayment) ProcessPayment(payment *models.Payment) error {
	fmt.Println("kaspi")
	return nil
}
func (p *KaspiPayment) CheckPaymentStatus(paymentID string) string {
	return ""
}

type PayPalPayment struct{}

func NewPayPalPayment() *PayPalPayment {
	return &PayPalPayment{}
}

func (p *PayPalPayment) ProcessPayment(payment *models.Payment) error {
	fmt.Println("paypal")
	return nil
}
func (p *PayPalPayment) CheckPaymentStatus(paymentID string) string {
	return ""
}
