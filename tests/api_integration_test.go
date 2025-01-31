//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/go-rest-user/config"
	"github.com/yakob-abada/go-rest-user/pkg/bootstrap"
	"github.com/yakob-abada/go-rest-user/pkg/migration"
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var e *echo.Echo

// Setup and teardown for integration tests
func TestMain(m *testing.M) {
	// Load .env.test to use test database
	err := godotenv.Load("../.env.test")

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

	// Setup Echo instance and routes manually
	e = echo.New()

	rabbitMQ, _ := config.InitRabbitMQ()
	defer rabbitMQ.Close()

	userHandler := bootstrap.NewUserHandler(testDB, rabbitMQ, nil)

	e.POST("/users", userHandler.SaveUser)
	e.GET("/users", userHandler.GetUsers)
	e.DELETE("/users/:id", userHandler.DeleteUser)
	e.PUT("/users/:id", userHandler.UpdateUser)

	// Run tests
	code := m.Run()

	// Cleanup database after tests
	testDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	os.Exit(code)
}

// Test Get Users API (Pagination & Filters)
func TestGetUsers(t *testing.T) {
	// Insert test users
	users := []model.User{
		{FirstName: "Alice", LastName: "Smith", Email: "alice@example.com", Password: "securepass", Country: "UK"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", Password: "securepass", Country: "USA"},
	}
	testDB.Create(&users)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.GreaterOrEqual(t, int(response["total_users"].(float64)), 2, "Should return at least 2 users")
}

// Test Create User API
func TestCreateUser(t *testing.T) {
	userData := map[string]string{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "john.doe@example.com",
		"password":   "securepassword",
		"country":    "USA",
	}

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

// TestUpdateUser tests updating an existing user
func TestUpdateUser(t *testing.T) {
	// Create test user
	user := model.User{
		FirstName: "Mark",
		LastName:  "Taylor",
		Email:     "mark.taylor@example.com",
		Password:  "securehash",
		Country:   "Canada",
	}
	testDB.Create(&user)

	// Make UPDATE request
	updateData := map[string]interface{}{
		"first_name": "Johnny",
		"country":    "Canada",
	}

	jsonData, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPut, "/users/"+user.ID.String(), bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var updatedUser model.User
	err := json.Unmarshal(rec.Body.Bytes(), &updatedUser)
	assert.NoError(t, err)

	assert.Equal(t, "Johnny", updatedUser.FirstName)
	assert.Equal(t, "Canada", updatedUser.Country)
	assert.Equal(t, user.Email, updatedUser.Email) // Email should remain unchanged
	assert.Equal(t, http.StatusOK, rec.Code)
}

// Test Delete User API
func TestDeleteUser(t *testing.T) {
	// Create test user
	user := model.User{
		FirstName: "Mark",
		LastName:  "Taylor",
		Email:     "mark.taylor@example.com",
		Password:  "securehash",
		Country:   "Canada",
	}
	testDB.Create(&user)

	// Make DELETE request
	req := httptest.NewRequest(http.MethodDelete, "/users/"+user.ID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
