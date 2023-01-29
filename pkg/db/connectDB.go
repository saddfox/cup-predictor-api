package db

import (
	"fmt"
	"log"
	"os"

	"github.com/saddfox/cup-predictor/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// initialize connection to the database and migrate models
func ConnectDB() {
	var err error
	dbURL := fmt.Sprintf("postgres://postgres:postgrespwd@%s/db", os.Getenv("POSTGRES_URL"))

	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to DB")

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Format1{})
	DB.AutoMigrate(&models.Format2{})
	DB.AutoMigrate(&models.Cup{})
	fmt.Println("Migrated DB schema")
}
