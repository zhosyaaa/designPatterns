package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pattern/internal/repositories"
	"pattern/internal/services"
)

type TransactionController struct {
	transactionFactories map[string]services.TransactionFactory
}

func NewTransactionController() *TransactionController {
	controller := &TransactionController{
		transactionFactories: make(map[string]services.TransactionFactory),
	}
	controller.RegisterTransactionFactory("transfer", &repositories.TransferTransaction{})
	controller.RegisterTransactionFactory("withdrawals", &repositories.WithdrawalsTransaction{})
	controller.RegisterTransactionFactory("deposit", &repositories.DepositTransaction{})
	return controller
}

func (c *TransactionController) RegisterTransactionFactory(transactionType string, factory services.TransactionFactory) {
	c.transactionFactories[transactionType] = factory
}

func (c *TransactionController) Transaction(context *gin.Context) {
	transactionType := context.Param("type")
	factory, exists := c.transactionFactories[transactionType]
	if !exists {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported transaction type"})
		return
	}

	transaction := factory.CreateTransaction()
	transaction.Execute()

	context.JSON(http.StatusOK, gin.H{"message": "Transaction executed successfully"})
}
