package errorhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/go-rest-user/pkg/common"
	"github.com/yakob-abada/go-rest-user/pkg/logging"
)

// TestHandleBadRequest ensures bad request errors are logged and return the correct response
func TestHandleBadRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLogger := new(logging.MockLogger)
	errorHandler := NewErrorHandler(mockLogger)

	// Mock logger expectation
	mockLogger.On("Warn", mock.Anything, "Invalid input", mock.Anything).Return()

	// Context with Correlation ID
	ctx := context.WithValue(context.Background(), common.CorrelationID, "test-correlation-id")

	details := map[string]interface{}{
		"field": "email",
		"error": "invalid format",
	}

	err := errorHandler.HandleBadRequest(ctx, c, "Invalid input", details)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, "Invalid input", response["error"])
	assert.Equal(t, details, response["details"])
	assert.Equal(t, "test-correlation-id", response["correlation_id"])

	mockLogger.AssertExpectations(t)
}

// TestHandleInternalServerError ensures internal server errors are logged and return the correct response
func TestHandleInternalServerError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLogger := new(logging.MockLogger)
	errorHandler := NewErrorHandler(mockLogger)

	// Mock logger expectation
	mockLogger.On("Error", mock.Anything, "Database error", mock.Anything).Return()

	// Context with Correlation ID
	ctx := context.WithValue(context.Background(), common.CorrelationID, "test-correlation-id")

	details := map[string]interface{}{
		"query": "SELECT * FROM users",
		"error": "timeout",
	}

	err := errorHandler.HandleInternalServerError(ctx, c, "Database error", details)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, "Database error", response["error"])
	assert.Equal(t, details, response["details"])
	assert.Equal(t, "test-correlation-id", response["correlation_id"])

	mockLogger.AssertExpectations(t)
}
