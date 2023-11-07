package services

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
	"net/smtp"
	"pattern/internal/clients"
	"pattern/internal/models"
	"pattern/internal/repositories"
	"time"
)

type CurrencyService struct {
	userRepo repositories.UserRepository
	client   clients.CurrencyClient
}

func NewCurrencyService(userRepo repositories.UserRepository) *CurrencyService {
	return &CurrencyService{userRepo: userRepo}
}

func (c CurrencyService) Notify(user *models.User, currency string, newRate float64) error {
	subject := "Currency exchange rate change"
	body := fmt.Sprintf("The exchange rate of %s has changed by %.2f", currency, newRate)

	err := c.SendEmail(user.Email, subject, body)
	if err != nil {
		return err
	}

	return nil
}

func (c CurrencyService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "SMTP for My App (musabecova05@gmail.com)")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	buf := new(bytes.Buffer)
	m.WriteTo(buf)

	appPassword := "mayf ayum loqn haqs"
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("gmail", "musabecova05@gmail.com", appPassword, "smtp.gmail.com"), "musabecova05@gmail.com", []string{to}, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c CurrencyService) CheckUpdates() {
	lastRates := make(map[string]float64)

	for {
		newRates, err := c.client.GetExchangeRates()
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch currency rates")
			time.Sleep(1 * time.Minute)
			continue
		}

		for currency, newRate := range newRates {
			lastRate, ok := lastRates[currency]
			if !ok {
				lastRates[currency] = newRate
				continue
			}

			if newRate != lastRate {
				subscribers, err := c.userRepo.GetSubscribersByCurrency(currency)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get subscribers for currency %s", currency)
					continue
				}

				for _, subscriber := range subscribers {
					err := c.Notify(&subscriber, currency, newRate)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to send notification to user %s", subscriber.Username)
					}
				}

				lastRates[currency] = newRate
			}
		}
		time.Sleep(5 * time.Minute)
		lastRates["KZT"] = 501

	}
}
