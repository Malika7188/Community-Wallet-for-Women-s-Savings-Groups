// database/db.go
package database

import (
	"fmt"
	"log"
	"os"

	// "your_project_name/smodels" // replace with your actual module

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"chama-wallet-backend/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.Group{}, &models.Member{}, &models.Contribution{})
	err = DB.AutoMigrate(&models.Group{}, &models.Member{}, &models.Contribution{}, &models.User{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	fmt.Println("Connected to the database successfully.")
}
