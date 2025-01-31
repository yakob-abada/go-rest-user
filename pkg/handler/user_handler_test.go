//go:build unit_test

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/go-rest-user/pkg/errorhandler"
	"github.com/yakob-abada/go-rest-user/pkg/logging"
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"github.com/yakob-abada/go-rest-user/pkg/publisher"
	"github.com/yakob-abada/go-rest-user/pkg/repository"
	"github.com/yakob-abada/go-rest-user/pkg/security"
	"github.com/yakob-abada/go-rest-user/pkg/validator"
)

func TestGetUsers_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users?page=1&limit=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	users := []model.User{
		{ID: uuid.New(), FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Country: "USA"},
		{ID: uuid.New(), FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com", Country: "Canada"},
	}

	// Define expectations
	mockRepo.On("GetUsers", mock.Anything, 1, 2, mock.Anything, "", "").Return(users, int64(2), nil)

	handler := NewUserHandler(mockRepo, mockLogger, nil, nil, mockErrorHandler, mockHasher)

	err := handler.GetUsers(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), resp["total"])

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestSaveUser_Success(t *testing.T) {
	e := echo.New()
	reqBody := `{"first_name":"John","last_name":"Doe","email":"john.doe@example.com","password":"securepass","country":"USA"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockPublisher := new(publisher.MockPublisher)
	mockValidator := new(validator.MockValidator)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	// Define expectations
	mockValidator.On("ValidateUser", mock.Anything).Return(nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").Return(nil, nil)
	mockHasher.On("HashPassword", "securepass").Return("$2a$10$hashedpassword123", nil)
	mockRepo.On("SaveUser", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		user := args.Get(1).(*model.User)
		user.ID = uuid.New()
	}).Return(nil)
	mockPublisher.On("Publish", publisher.EventUserCreated, mock.Anything).Return(nil)

	handler := NewUserHandler(mockRepo, mockLogger, mockValidator, mockPublisher, mockErrorHandler, mockHasher)

	err := handler.SaveUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp model.User
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "John", resp.FirstName)
	assert.NotEqual(t, "securepass", resp.Password)
}

func TestSaveUser_DuplicateEmail(t *testing.T) {
	e := echo.New()
	reqBody := `{"first_name":"John","last_name":"Doe","email":"john.doe@example.com","password":"securepass","country":"USA"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockPublisher := new(publisher.MockPublisher)
	mockValidator := new(validator.MockValidator)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	existingUser := &model.User{Email: "john.doe@example.com"}
	mockValidator.On("ValidateUser", mock.Anything).Return(nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").Return(existingUser, nil)

	// Expect logger to log duplication attempt
	mockLogger.On("Warn", mock.Anything, "User already exists", mock.Anything).Return()

	mockErrorHandler.On("HandleBadRequest", mock.Anything, c, "User with this email already exists", mock.Anything).
		Return(c.JSON(http.StatusConflict, map[string]string{"error": "User with this email already exists"}))

	handler := NewUserHandler(mockRepo, mockLogger, mockValidator, mockPublisher, mockErrorHandler, mockHasher)

	err := handler.SaveUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "User with this email already exists")

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}

func TestSaveUser_ValidationFailure(t *testing.T) {
	e := echo.New()
	reqBody := `{"email":"invalid-email"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockPublisher := new(publisher.MockPublisher)
	mockValidator := new(validator.MockValidator)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	// Define expectations
	mockValidator.On("ValidateUser", mock.Anything).Return(assert.AnError) // Validation fails
	mockErrorHandler.On("HandleBadRequest", mock.Anything, c, "Validation failed", mock.Anything).
		Return(c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"}))

	handler := NewUserHandler(mockRepo, mockLogger, mockValidator, mockPublisher, mockErrorHandler, mockHasher)

	err := handler.SaveUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Validation failed")

	// Verify mock expectations
	mockValidator.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}

// Test UpdateUser Handler
func TestUpdateUser(t *testing.T) {
	// Initialize Echo instance and Mock Repository
	e := echo.New()

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockPublisher := new(publisher.MockPublisher)
	mockValidator := new(validator.MockValidator)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	// Sample User Update Data
	userID := "550e8400-e29b-41d4-a716-446655440000"

	updateData := &model.UserUpdate{
		FirstName: "Johnny",
		Country:   "Canada",
	}

	expectedUser := &model.User{
		ID:        uuid.MustParse(userID),
		FirstName: "Johnny",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Country:   "Canada",
	}

	mockRepo.On("UpdateUser", mock.Anything, userID, updateData).Return(expectedUser, nil)
	mockPublisher.On("Publish", publisher.EventUserUpdated, mock.Anything).Return(nil)

	// Convert updateData to JSON
	jsonData, _ := json.Marshal(updateData)

	// Create a request
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID+"?todo=1", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Set Echo Context
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	// Call Handler
	handler := NewUserHandler(mockRepo, mockLogger, mockValidator, mockPublisher, mockErrorHandler, mockHasher)
	err := handler.UpdateUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse Response Body
	var response model.User
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate Response Fields
	assert.Equal(t, expectedUser.FirstName, response.FirstName)
	assert.Equal(t, expectedUser.Country, response.Country)

	mockRepo.AssertExpectations(t)
}

// Test UpdateUser with Invalid JSON Body
func TestUpdateUser_InvalidBody(t *testing.T) {
	e := echo.New()

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockPublisher := new(publisher.MockPublisher)
	mockValidator := new(validator.MockValidator)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	// Create Request with Invalid JSON
	req := httptest.NewRequest(http.MethodPut, "/users/123", bytes.NewReader([]byte("{invalid_json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Set Echo Context
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	// Call Handler
	handler := NewUserHandler(mockRepo, mockLogger, mockValidator, mockPublisher, mockErrorHandler, mockHasher)
	err := handler.UpdateUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Validate Error Response
	expectedResponse := `{"error":"Invalid request body"}`
	assert.JSONEq(t, expectedResponse, rec.Body.String())

	mockRepo.AssertExpectations(t)
}

// Test UpdateUser When User Not Found
func TestUpdateUser_UserNotFound(t *testing.T) {
	e := echo.New()

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockPublisher := new(publisher.MockPublisher)
	mockValidator := new(validator.MockValidator)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	userID := "550e8400-e29b-41d4-a716-446655440000"

	updateData := &model.UserUpdate{
		FirstName: "Johnny",
	}

	mockRepo.On("UpdateUser", mock.Anything, userID, updateData).Return(nil, errors.New("user not found"))

	// Convert updateData to JSON
	jsonData, _ := json.Marshal(updateData)

	// Create Request
	req := httptest.NewRequest(http.MethodPut, "/users/"+userID, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Set Echo Context
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID)

	// Call Handler
	handler := NewUserHandler(mockRepo, mockLogger, mockValidator, mockPublisher, mockErrorHandler, mockHasher)
	err := handler.UpdateUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Validate Error Response
	expectedResponse := `{"error":"user not found"}`
	assert.JSONEq(t, expectedResponse, rec.Body.String())

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	// Mock dependencies
	mockRepo := new(repository.MockUserRepository)
	mockPublisher := new(publisher.MockPublisher)
	mockLogger := new(logging.MockLogger)
	mockErrorHandler := new(errorhandler.MockErrorHandler)
	mockHasher := new(security.MockPasswordHasher)

	// Define expectations
	mockRepo.On("DeleteUser", mock.Anything, userID).Return(nil)
	mockPublisher.On("Publish", publisher.EventUserDeleted, mock.Anything).Return(nil)

	handler := NewUserHandler(mockRepo, mockLogger, nil, mockPublisher, mockErrorHandler, mockHasher)

	err := handler.DeleteUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "User deleted successfully")

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
