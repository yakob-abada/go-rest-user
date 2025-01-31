package publisher

import (
	"github.com/stretchr/testify/mock"
)

// MockPublisher implements the Publisher interface for testing
type MockPublisher struct {
	mock.Mock
}

// Publish mocks the Publish method
func (m *MockPublisher) Publish(event EventType, payload map[string]interface{}) error {
	args := m.Called(event, payload)
	return args.Error(0)
}
