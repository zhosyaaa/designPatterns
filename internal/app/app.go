package app

import (
	"log"
	"pattern/internal/controllers"
	"pattern/internal/db"
	"pattern/internal/helpers"
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
	router := routes.NewRoutes(*pay, *auth)
	counterSource := helpers.NewCounterMetricsSource()
	timerSource := helpers.NewTimerMetricsSource()

	prometheusAdapter := helpers.NewPrometheusAdapter()

	counterSource.RegisterCounter("requests_total")
	counterSource.IncrementCounter("requests_total")

	timerSource.RegisterTimer("request_duration_seconds")
	timerSource.RecordTime("request_duration_seconds", 0.5)

	prometheusAdapter.RegisterCounter("prometheus_requests_total")
	prometheusAdapter.IncrementCounter("prometheus_requests_total")

	prometheusAdapter.RegisterTimer("prometheus_request_duration_seconds")
	prometheusAdapter.RecordTime("prometheus_request_duration_seconds", 1.0)

	return *router
}
