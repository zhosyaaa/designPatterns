package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"pattern/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	fmt.Println(r.db)
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) AddPaymentToUser(userID uint, payment *models.Payment) error {
	user, err := u.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.Payments = append(user.Payments, *payment)
	result := u.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserRepository) CreateSubscription(subscription *models.Subscription) error {
	if err := u.db.Create(subscription).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) DeleteSubscription(userID uint, currency string) error {
	if err := u.db.Where("user_id = ? AND currency = ?", userID, currency).Delete(&models.Subscription{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetSubscribersByCurrency(currency string) ([]models.User, error) {
	var users []models.User

	if err := r.db.Table("subscriptions").
		Select("users.*").
		Joins("JOIN users ON subscriptions.user_id = users.id").
		Where("subscriptions.currency = ? AND subscriptions.deleted_at IS NULL", currency).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
