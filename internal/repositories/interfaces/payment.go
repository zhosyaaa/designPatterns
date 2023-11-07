package interfaces

import "pattern/internal/models"

type ProcessPayment interface {
	ProcessPayment(payment *models.Payment) error
}

type PaymentService interface {
	CreatePayment(payment *models.Payment) (string, error)
}

type CurrencyConverter interface {
	Convert(amount int64, fromCurrency string, toCurrency string) (float64, error)
}
type CurrencyRateProvider interface {
	GetCurrencyRates() (map[string]float64, error)
}
