package mapper

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/go-rest-user/pkg/model"
)

// TestUserResponseMapper ensures correct mapping from User to UserResponse.
func TestUserResponseMapper(t *testing.T) {
	userID := uuid.New() // Generate a new UUID

	user := model.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "JD",
		Email:     "john.doe@example.com",
		Country:   "USA",
	}

	expectedResponse := model.UserResponse{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "JD",
		Email:     "john.doe@example.com",
		Country:   "USA",
	}

	response := UserResponseMapper(user)

	assert.Equal(t, expectedResponse, response, "UserResponseMapper should correctly map User to UserResponse")
}

// TestUsersResponseMapper ensures correct mapping from []User to []UserResponse.
func TestUsersResponseMapper(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()

	users := []model.User{
		{
			ID:        userID1,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "JD",
			Email:     "john.doe@example.com",
			Country:   "USA",
		},
		{
			ID:        userID2,
			FirstName: "Jane",
			LastName:  "Smith",
			Nickname:  "JS",
			Email:     "jane.smith@example.com",
			Country:   "Canada",
		},
	}

	expectedResponses := []model.UserResponse{
		{
			ID:        userID1,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "JD",
			Email:     "john.doe@example.com",
			Country:   "USA",
		},
		{
			ID:        userID2,
			FirstName: "Jane",
			LastName:  "Smith",
			Nickname:  "JS",
			Email:     "jane.smith@example.com",
			Country:   "Canada",
		},
	}

	responses := UsersResponseMapper(users)

	assert.Equal(t, expectedResponses, responses, "UsersResponseMapper should correctly map []User to []UserResponse")
}
