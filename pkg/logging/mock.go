package logging

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	mock.Mock
}

// Info mocks the Info logging method
func (m *MockLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	m.Called(ctx, msg, fields)
}

// Warn mocks the Warn logging method
func (m *MockLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	m.Called(ctx, msg, fields)
}

// Error mocks the Error logging method
func (m *MockLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	m.Called(ctx, msg, fields)
}
