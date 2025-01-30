package security

import (
	"github.com/stretchr/testify/mock"
)

// MockPasswordHasher implements the PasswordHasher interface for testing
type MockPasswordHasher struct {
	mock.Mock
}

// HashPassword mocks the hashing method
func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// ComparePassword mocks the password comparison
func (m *MockPasswordHasher) ComparePassword(hashedPassword, plainPassword string) bool {
	args := m.Called(hashedPassword, plainPassword)
	return args.Bool(0)
}
