package app

import (
	"log"
	"pattern/internal/clients"
	"pattern/internal/controllers"
	"pattern/internal/db"
	"pattern/internal/middlewares"
	"pattern/internal/repositories"
	"pattern/internal/repositories/interfaces"
	"pattern/internal/routes"
	"pattern/internal/services"
)

func SetupApp() routes.Routes {
	dbInstance, err := db.GetDBInstance()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	userRepo := repositories.NewUserRepository(dbInstance)
	auth := controllers.NewAuthController(*userRepo)
	var curConv interfaces.CurrencyConverter
	kaspi := &middlewares.CurrencyConversionDecorator{&middlewares.BalanceCheckingDecorator{&repositories.KaspiPayment{}, *userRepo}, curConv}
	payPal := &middlewares.CurrencyConversionDecorator{&middlewares.BalanceCheckingDecorator{&repositories.KaspiPayment{}, *userRepo}, curConv}
	purchase := services.NewPurchase(dbInstance)
	purchase.RegisterStrategy("Kaspi", kaspi)
	purchase.RegisterStrategy("PayPal", payPal)
	pay := controllers.NewPaymentController(*purchase, *userRepo)
	client := &clients.CurrencyClient{}
	currencyService := services.NewCurrencyService(*userRepo)
	currency := controllers.NewCurrencyController(*userRepo, *client, *currencyService)
	go currencyService.CheckUpdates()
	router := routes.NewRoutes(*pay, *auth, *currency)
	return *router
}
