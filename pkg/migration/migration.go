package migration

import (
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"gorm.io/gorm"
	"log"
)

// RunMigrations applies database migrations
func RunMigrations(db *gorm.DB) {
	log.Println("🚀 Running database migrations...")

	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
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
