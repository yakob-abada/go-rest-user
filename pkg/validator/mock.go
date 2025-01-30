package validator

import (
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/go-rest-user/pkg/model"
)

// MockValidator implements the Validator interface for testing
type MockValidator struct {
	mock.Mock
}

// ValidateUser mocks the ValidateUser method
func (m *MockValidator) ValidateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
