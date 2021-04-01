package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	SenderID   uint64
	Sender     User
	ReceiverID uint64
	Receiver   User
	Amount     uint64
}
