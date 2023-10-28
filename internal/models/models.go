package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Payments []Payment `json:"payments"`
	Balance  int64     `json:"balance"`
	Currency string    `json:"currency"`
	Password string
}

type Payment struct {
	gorm.Model
	UserID   uint   `gorm:"not null" json:"-"`
	User     User   `gorm:"foreignkey:UserID"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Subscription struct {
	ID            uint
	UserID        uint
	User          User `gorm:"foreignkey:UserID"`
	Currency      string
	NotifyAddress string
}
