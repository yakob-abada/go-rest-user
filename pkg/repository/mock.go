package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/go-rest-user/pkg/model"
)

// MockUserRepository implements the UserRepository interface for testing
type MockUserRepository struct {
	mock.Mock
}

// GetUsers mocks the GetUsers method
func (m *MockUserRepository) GetUsers(ctx context.Context, page, limit int, filters map[string]string, sortBy, order string) ([]model.User, int64, error) {
	args := m.Called(ctx, page, limit, filters, sortBy, order)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

// GetUserByID mocks the GetUserByID method
func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetUserByEmail mocks the GetUserByEmail method
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserRepository) UpdateUser(ctx context.Context, id string, updatedData *model.UserUpdate) (*model.User, error) {
	args := m.Called(ctx, id, updatedData)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// SaveUser mocks the SaveUser method
func (m *MockUserRepository) SaveUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
