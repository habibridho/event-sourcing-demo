package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	SenderID   uint
	Sender     User
	ReceiverID uint
	Receiver   User
	Amount     uint64
}
