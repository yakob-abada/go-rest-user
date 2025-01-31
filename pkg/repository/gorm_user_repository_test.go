//go:build integration

package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/go-rest-user/pkg/migration"
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"gorm.io/driver/postgres"
)

var testDB *gorm.DB

// TestMain sets up the test database before running tests
func TestMain(m *testing.M) {
	// Load .env.test file
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Println("⚠️ Warning: No .env.test file found, using default environment variables")
	}

	// Read database environment variables
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	// Construct the database connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	// Connect to PostgreSQL
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to test database: %v", err)
	}

	log.Println("✅ Connected to test database!")

	// Run migrations before tests
	migration.RunMigrations(testDB)

	// Run tests
	code := m.Run()

	// Cleanup: Reset database after tests
	log.Println("🧹 Cleaning up test database...")
	testDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	os.Exit(code)
}

// cleanTestDatabase resets all tables before each test
func cleanTestDatabase() {
	log.Println("🧹 Truncating tables before test execution...")
	testDB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;")
}

func TestGetUsers(t *testing.T) {
	cleanTestDatabase()
	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	users := []model.User{
		{ID: uuid.New(), FirstName: "Mike", LastName: "Tyson", Email: "mike@example.com", Password: "pass123", Country: "USA"},
		{ID: uuid.New(), FirstName: "Serena", LastName: "Williams", Email: "serena@example.com", Password: "pass123", Country: "USA"},
	}

	for _, u := range users {
		assert.NoError(t, repo.SaveUser(ctx, &u))
	}

	// Fetch users
	fetchedUsers, total, err := repo.GetUsers(ctx, 1, 10, map[string]string{"country": "USA"}, "first_name", "asc")
	assert.NoError(t, err)
	assert.Len(t, fetchedUsers, 2)
	assert.Equal(t, int64(2), total)
}
func TestGetUserByEmail(t *testing.T) {
	cleanTestDatabase()
	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	user := model.User{
		ID:        uuid.New(),
		FirstName: "Alice",
		LastName:  "Wonderland",
		Email:     "alice@example.com",
		Password:  "hashedpassword123",
		Country:   "UK",
	}

	// Save user
	err := repo.SaveUser(ctx, &user)
	assert.NoError(t, err)

	// Fetch existing user
	storedUser, err := repo.GetUserByEmail(ctx, "alice@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, storedUser)
	assert.Equal(t, "Alice", storedUser.FirstName)

	// Fetch non-existing user
	nonExistentUser, err := repo.GetUserByEmail(ctx, "unknown@example.com")
	assert.NoError(t, err)
	assert.Nil(t, nonExistentUser)
}

func TestSaveUser(t *testing.T) {
	cleanTestDatabase()
	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	user := model.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "hashedpassword123",
		Country:   "USA",
	}

	// Save user
	err := repo.SaveUser(ctx, &user)
	assert.NoError(t, err)

	// Retrieve user by email
	storedUser, err := repo.GetUserByEmail(ctx, "john.doe@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, storedUser)
	assert.Equal(t, "John", storedUser.FirstName)
	assert.Equal(t, "Doe", storedUser.LastName)
	assert.Equal(t, "john.doe@example.com", storedUser.Email)
}

func TestUpdateUser(t *testing.T) {
	cleanTestDatabase()
	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	// Create a test user
	user := model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "hashedpassword",
		Country:   "USA",
	}
	// Save user
	err := repo.SaveUser(ctx, &user)
	assert.NoError(t, err)

	// Update user
	updatedData := &model.UserUpdate{
		FirstName: "Johnny",
		Country:   "Canada",
	}

	updatedUser, err := repo.UpdateUser(ctx, user.ID.String(), updatedData)
	assert.NoError(t, err)
	assert.Equal(t, "Johnny", updatedUser.FirstName)
	assert.Equal(t, "Canada", updatedUser.Country)
}

func TestDeleteUser(t *testing.T) {
	cleanTestDatabase()
	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	user := model.User{
		ID:        uuid.New(),
		FirstName: "Bob",
		LastName:  "Marley",
		Email:     "bob@example.com",
		Password:  "securepass",
		Country:   "Jamaica",
	}

	// Save user
	err := repo.SaveUser(ctx, &user)
	assert.NoError(t, err)

	// Delete user
	err = repo.DeleteUser(ctx, user.ID)
	assert.NoError(t, err)

	// Verify user is deleted
	deletedUser, err := repo.GetUserByEmail(ctx, "bob@example.com")
	assert.NoError(t, err)
	assert.Nil(t, deletedUser)
}
