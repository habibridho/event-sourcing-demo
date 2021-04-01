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

	log.Print("seeding db with sample data")
	seedSampleData()

	log.Print("db connection initialised.")
}

func seedSampleData() {
	sampleUser1 := model.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email: "habib@email.com",
	}
	sampleUser2 := model.User{
		Model: gorm.Model{
			ID: 2,
		},
		Email: "ani@email.com",
	}
	sampleUser3 := model.User{
		Model: gorm.Model{
			ID: 3,
		},
		Email: "budi@email.com",
	}
	sampleUser1.SetPassword("password")
	sampleUser2.SetPassword("password")
	sampleUser3.SetPassword("password")
	db.Create(&sampleUser1)
	db.Create(&sampleUser2)
	db.Create(&sampleUser3)

	sampleAccount := model.Account{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:  1,
		Balance: 9000000000,
	}
	sampleAccount2 := model.Account{
		Model: gorm.Model{
			ID: 2,
		},
		UserID:  2,
		Balance: 0,
	}
	sampleAccount3 := model.Account{
		Model: gorm.Model{
			ID: 3,
		},
		UserID:  3,
		Balance: 0,
	}
	db.Create(&sampleAccount)
	db.Create(&sampleAccount2)
	db.Create(&sampleAccount3)
}
