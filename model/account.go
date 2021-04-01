package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID  uint64
	User    User
	Balance uint64
}
