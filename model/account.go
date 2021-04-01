package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID  uint
	User    User
	Balance uint64
}
