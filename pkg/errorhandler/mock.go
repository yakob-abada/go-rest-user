package errorhandler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockErrorHandler implements the ErrorHandler interface for testing
type MockErrorHandler struct {
	mock.Mock
}

// HandleBadRequest mocks the HandleBadRequest method
func (m *MockErrorHandler) HandleBadRequest(ctx context.Context, c echo.Context, message string, details map[string]interface{}) error {
	args := m.Called(ctx, c, message, details)
	return args.Error(0)
}

// HandleInternalServerError mocks the HandleInternalServerError method
func (m *MockErrorHandler) HandleInternalServerError(ctx context.Context, c echo.Context, message string, details map[string]interface{}) error {
	args := m.Called(ctx, c, message, details)
	return args.Error(0)
}

// HandleNotFound mocks the HandleNotFound method
func (m *MockErrorHandler) HandleNotFound(ctx context.Context, c echo.Context, message string, details map[string]interface{}) error {
	args := m.Called(ctx, c, message, details)
	return args.Error(0)
}
