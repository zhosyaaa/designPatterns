package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"pattern/internal/clients"
	"pattern/internal/models"
	"pattern/internal/repositories"
	"pattern/internal/services"
)

type CurrencyController struct {
	userRepo        repositories.UserRepository
	client          clients.CurrencyClient
	currencyService services.CurrencyService
}

func NewCurrencyController(userRepo repositories.UserRepository, client clients.CurrencyClient, currencyService services.CurrencyService) *CurrencyController {
	return &CurrencyController{userRepo: userRepo, client: client, currencyService: currencyService}
}

func (c CurrencyController) GetCurrencies(context *gin.Context) {
	symbol := context.Param("symbol")

	if symbol == "" {
		log.Error().Msg("Symbol not provided")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not provided"})
		return
	}
	adapter := services.CurrencyClientAdapter{Client: &c.client}
	currencyData, err := adapter.GetCurrencyRates()
	if err != nil {
		log.Error().Msg("Symbol not provided")
		context.JSON(http.StatusBadRequest, gin.H{"error": "rates not found"})
		return
	}

	for _, data := range currencyData {
		if data.Name == symbol {
			log.Info().Msgf("Currency found: Symbol=%s, Price=%.2f", data.Name, data.Value)
			context.JSON(http.StatusOK, gin.H{"symbol": data.Name, "price": data.Value})
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

	if err = c.currencyService.Notify(user, "KZT", 500); err != nil {
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
