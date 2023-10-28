package interfaces

import "pattern/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateSubscription(subscription *models.Subscription) error
	DeleteSubscription(userID uint, currency string) error
	GetSubscribersByCurrency(currency string) ([]models.User, error)
}
