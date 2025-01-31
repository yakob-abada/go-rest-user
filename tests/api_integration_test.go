//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/go-rest-user/config"
	"github.com/yakob-abada/go-rest-user/pkg/bootstrap"
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var e *echo.Echo

// Setup and teardown for integration tests
func TestMain(m *testing.M) {
	// Load .env.test to use test database
	_ = godotenv.Load("../.env.test")

	// Initialize test database
	testDB := config.InitDB()

	// Run migrations for testing
	_ = testDB.AutoMigrate(&model.User{})

	// Setup Echo instance and routes manually
	e = echo.New()
	userHandler := bootstrap.NewUserHandler(testDB, nil, nil)

	e.POST("/users", userHandler.SaveUser)
	e.GET("/users", userHandler.GetUsers)
	e.DELETE("/users/:id", userHandler.DeleteUser)

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

	assert.GreaterOrEqual(t, int(response["total"].(float64)), 2, "Should return at least 2 users")
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

// Test Get User by ID API
func TestGetUserByID(t *testing.T) {
	// Create test user
	user := model.User{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice.smith@example.com",
		Password:  "hashedpassword",
		Country:   "UK",
	}
	testDB.Create(&user)

	// Make GET request
	req := httptest.NewRequest(http.MethodGet, "/users/"+user.ID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

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
