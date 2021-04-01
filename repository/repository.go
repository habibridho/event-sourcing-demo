package repository

import (
	"context"
	"errors"
	"event-sourcing-demo/model"
	"gorm.io/gorm"
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
