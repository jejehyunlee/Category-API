package database

import (
	"Category-API/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Prioritize DATABASE_URL from environment
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		// Fallback to individual variables
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbSSLMode := os.Getenv("DB_SSLMODE")

		if dbSSLMode == "" {
			dbSSLMode = "disable"
		}

		databaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
	}

	log.Printf("Connecting to database...")

	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")

	// Auto migrate
	err = DB.AutoMigrate(&models.Category{})
	if err != nil {
		log.Printf("Warning: AutoMigrate failed: %v", err)
	} else {
		log.Println("Database migration completed")
	}

	// Test connection
	testConnection()
}

func testConnection() {
	var result int
	if err := DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		log.Printf("Database test query failed: %v", err)
	} else {
		log.Printf("Database test query successful: %d", result)
	}
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
