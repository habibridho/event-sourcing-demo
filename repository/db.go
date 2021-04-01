package repository

import (
	"event-sourcing-demo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func InitialiseDB() {
	dsn := "host=localhost user=postgres password=password dbname=emoney port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	log.Print("migrating model...")
	db.AutoMigrate(&model.User{},
		&model.Account{},
		&model.Transaction{})
	log.Print("model migrated.")
	log.Print("db connection initialised.")
}
