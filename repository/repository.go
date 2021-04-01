package repository

import (
	"context"
	"errors"
	"event-sourcing-demo/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

func FetchUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	if err := db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return user, err
		}
	}
	return user, nil
}

func ExecuteTransaction(ctx context.Context, transaction model.Transaction) (err error) {
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	senderID := transaction.SenderID
	receiverID := transaction.ReceiverID
	amount := transaction.Amount

	var senderAccount, receiverAccount model.Account
	if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&senderAccount, "user_id = ?", senderID).Error; err != nil {
		log.Printf("could not fetch sender account: %s", err.Error())
		return
	}
	if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&receiverAccount, "user_id = ?", receiverID).Error; err != nil {
		log.Printf("could not fetch receiver account: %s", err.Error())
		return
	}

	if senderAccount.Balance < amount {
		err = InsufficientBalance{}
		return
	}
	senderAccount.Balance -= amount
	receiverAccount.Balance += amount

	if err = tx.Create(&transaction).Error; err != nil {
		log.Printf("could not save transaction data: %s", err.Error())
		return
	}
	if err = tx.Save(&senderAccount).Error; err != nil {
		log.Printf("could not update sender account: %s", err.Error())
		return
	}
	if err = tx.Save(&receiverAccount).Error; err != nil {
		log.Printf("could not update receiver account: %s", err.Error())
		return
	}
	return
}
