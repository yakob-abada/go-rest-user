//go:build unit_test

package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/go-rest-user/pkg/model"
)

func TestUserValidator_ValidateUser(t *testing.T) {
	validator := NewDefaultUserValidator()

	tests := []struct {
		name    string
		user    model.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid User",
			user: model.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Country:   "USA",
			},
			wantErr: false,
		},
		{
			name: "Missing First Name",
			user: model.User{
				LastName: "Doe",
				Email:    "john.doe@example.com",
				Password: "password123",
				Country:  "USA",
			},
			wantErr: true,
			errMsg:  "first_name is required",
		},
		{
			name: "Missing Last Name",
			user: model.User{
				FirstName: "John",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Country:   "USA",
			},
			wantErr: true,
			errMsg:  "last_name is required",
		},
		{
			name: "Invalid Email",
			user: model.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalid-email",
				Password:  "password123",
				Country:   "USA",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "Missing Email",
			user: model.User{
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password123",
				Country:   "USA",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "Missing Password",
			user: model.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Country:   "USA",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "Short Password",
			user: model.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "123",
				Country:   "USA",
			},
			wantErr: true,
			errMsg:  "password must be at least 6 characters long",
		},
		{
			name: "Missing Country",
			user: model.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
			},
			wantErr: true,
			errMsg:  "country is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUser(&tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
