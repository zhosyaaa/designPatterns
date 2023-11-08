package middlewares

import (
	"fmt"
	"pattern/internal/models"
	"pattern/internal/repositories/interfaces"
)

type CurrencyConversionDecorator struct {
	DecoratedService  interfaces.ProcessPayment
	CurrencyConverter interfaces.CurrencyConverter
}

//
//func (c *CurrencyConversionDecorator) Convert(amount int64, fromCurrency string, toCurrency string) (int64, error) {
//	client := clients.NewCurrencyClient()
//	exchangeRates, err := client.GetExchangeRates()
//	if err != nil {
//		return 0, err
//	}
//
//	//fromRate, existsFrom := exchangeRates[fromCurrency]
//	//toRate, existsTo := exchangeRates[toCurrency]
//	//
//	//if !existsFrom || !existsTo {
//	//	return 0, fmt.Errorf("Invalid currency conversion")
//	//}
//	//
//	////fromRateFloat, okFrom := fromRate.(int64)
//	////toRateFloat, okTo := toRate.(int64)
//	//
//	//if !okFrom || !okTo {
//	//	return 0, fmt.Errorf("Invalid exchange rate format")
//	//}
//	//
//	//convertedAmount := amount * (toRateFloat / fromRateFloat)
//	//return convertedAmount, nil
//}

func (c *CurrencyConversionDecorator) ProcessPayment(payment *models.Payment) error {
	fmt.Println("ProcessPayment start")
	fromCurrency := payment.User.Currency

	if fromCurrency != payment.Currency {
		//convertedAmount, err := c.Convert(payment.Amount, fromCurrency, payment.Currency)
		//if err != nil {
		//	return err
		//}
		convertedAmount := 123
		payment.Amount = int64(convertedAmount)
	}
	fmt.Println("ProcessPayment correct")
	return c.DecoratedService.ProcessPayment(payment)
}
