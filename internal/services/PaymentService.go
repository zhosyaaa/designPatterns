package services

import (
	"fmt"
	"gorm.io/gorm"
	"pattern/internal/models"
	"pattern/internal/repositories/interfaces"
)

type Purchase struct {
	paymentService map[string]interfaces.ProcessPayment
	db             *gorm.DB
}

func NewPurchase(db *gorm.DB) *Purchase {
	return &Purchase{paymentService: make(map[string]interfaces.ProcessPayment), db: db}
}

func (p *Purchase) RegisterStrategy(paymentMethod string, strategy interfaces.ProcessPayment) {
	p.paymentService[paymentMethod] = strategy
}

func (p *Purchase) ProcessPayment(payment *models.Payment, paymentMethod string) (string, error) {
	paymentMethodStrategy := p.paymentService[paymentMethod]
	if err := paymentMethodStrategy.ProcessPayment(payment); err != nil {
		return "", err
	}
	return "the process was successful", nil
}

func (s *Purchase) CreatePayment(payment *models.Payment) error {
	fmt.Println("userid: ", payment.UserID, payment.User.ID)
	if payment.Amount <= 0 {
		return fmt.Errorf("Invalid payment data")
	}
	result := s.db.Create(payment)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
