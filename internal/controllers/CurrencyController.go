package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
	"net/http"
	"net/smtp"
	"pattern/internal/clients"
	"pattern/internal/models"
	"pattern/internal/repositories"
	"time"
)

type CurrencyController struct {
	userRepo repositories.UserRepository
	client   clients.CurrencyClient
}

func NewCurrencyController(userRepo repositories.UserRepository, client clients.CurrencyClient) *CurrencyController {
	return &CurrencyController{userRepo: userRepo, client: client}
}

func (c CurrencyController) GetCurrencies(context *gin.Context) {
	symbol := context.Param("symbol")

	if symbol == "" {
		log.Error().Msg("Symbol not provided")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not provided"})
		return
	}
	currencyData, err := c.client.GetExchangeRates()
	if err != nil {
		log.Error().Msg("Symbol not provided")
		context.JSON(http.StatusBadRequest, gin.H{"error": "rates not found"})
		return
	}

	for key, value := range currencyData {
		if key == symbol {
			log.Info().Msgf("Currency found: Symbol=%s, Price=%.2f", key, value)
			context.JSON(http.StatusOK, gin.H{"symbol": key, "price": value})
			return
		}
	}

	log.Warn().Msgf("Currency not found: Symbol=%s", symbol)
	context.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cryptocurrency with symbol %s not found", symbol)})
}

func (c CurrencyController) SubscribeUser(context *gin.Context) {
	var subscriptionData struct {
		Currency string
	}

	if err := context.ShouldBindJSON(&subscriptionData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	userEmail, emailExists := context.Get("email")
	if !emailExists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}

	email, ok := userEmail.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert email to string"})
		return
	}

	user, err := c.userRepo.GetUserByEmail(email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	newSubscription := models.Subscription{
		UserID:        user.ID,
		Currency:      subscriptionData.Currency,
		NotifyAddress: email,
	}

	err = c.userRepo.CreateSubscription(&newSubscription)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Error when creating Subscription"})
		return
	}

	if err = c.Notify(user, "KZT", 500); err != nil {
		fmt.Println("notify", err)
	}

	log.Info().Msgf("User %s (ID: %d) subscribed to currency %s",
		user.Username, user.ID, newSubscription.Currency)
	context.JSON(http.StatusOK, gin.H{"message": "Subscription created successfully"})
}

func (c CurrencyController) UnsubscribeUser(context *gin.Context) {
	var unsubscribeData struct {
		Currency string
	}

	if err := context.ShouldBindJSON(&unsubscribeData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	userEmail, emailExists := context.Get("email")
	if !emailExists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}

	email, ok := userEmail.(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert email to string"})
		return
	}

	user, err := c.userRepo.GetUserByEmail(email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.userRepo.DeleteSubscription(user.ID, unsubscribeData.Currency); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}
	log.Info().Msgf("User %s (ID: %d) unsubscribed from currency %s",
		user.Username, user.ID, unsubscribeData.Currency)
	context.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

func (c CurrencyController) Notify(user *models.User, currency string, newRate float64) error {
	subject := "Currency exchange rate change"
	body := fmt.Sprintf("The exchange rate of %s has changed by %.2f", currency, newRate)

	err := c.SendEmail(user.Email, subject, body)
	if err != nil {
		return err
	}

	return nil
}

func (c CurrencyController) SendEmail(to, subject, body string) error {
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

func (c CurrencyController) CheckUpdates() {
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
		time.Sleep(1 * time.Minute)
		lastRates["KZT"] = 501

	}
}
