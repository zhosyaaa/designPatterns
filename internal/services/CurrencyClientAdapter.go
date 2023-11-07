package services

import (
	"pattern/internal/clients"
	"pattern/internal/models"
)

type CurrencyClientAdapter struct {
	Client *clients.CurrencyClient
}

func (adapter *CurrencyClientAdapter) GetCurrencyRates() ([]models.Currency, error) {
	rates, err := adapter.Client.GetExchangeRates()
	if err != nil {
		return nil, err
	}
	var arrRates []models.Currency

	for currencyCode, rate := range rates {
		currency := models.Currency{
			currencyCode,
			rate,
		}
		arrRates = append(arrRates, currency)
	}

	return arrRates, nil
}
