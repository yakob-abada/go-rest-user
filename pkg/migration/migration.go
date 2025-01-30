package migration

import (
	"fmt"
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Load environment variables from .env file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using default environment variables")
	}
}

// Connect to PostgreSQL
func connectDB() (*gorm.DB, error) {
	loadEnv()

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

// RunMigrations applies database migrations
func RunMigrations() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	log.Println("🚀 Running database migrations...")

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("❌ Failed to enable uuid-ossp: %v", err)
	}

	// Apply migrations for models
	err = db.AutoMigrate(&model.User{}) // Add more models as needed
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Database migration completed successfully!")
}
