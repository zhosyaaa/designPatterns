package middlewares

import (
	"fmt"
	"pattern/internal/models"
	"pattern/internal/repositories"
	"pattern/internal/repositories/interfaces"
)

type BalanceCheckingDecorator struct {
	DecoratedService interfaces.ProcessPayment
	UserRepository   repositories.UserRepository
}

func (b *BalanceCheckingDecorator) ProcessPayment(payment *models.Payment) error {
	fmt.Println(payment.User.Balance, payment.Amount)
	if payment.User.Balance < payment.Amount {
		return fmt.Errorf("Insufficient balance")
	}
	fmt.Println("BalanceCheckingDecorator correct")
	return b.DecoratedService.ProcessPayment(payment)
}
