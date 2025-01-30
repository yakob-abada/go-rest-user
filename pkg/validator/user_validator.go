package validator

import (
	"errors"
	"github.com/yakob-abada/go-rest-user/pkg/model"
	"net/mail"
)

// UserValidator defines an interface for user validation
type UserValidator interface {
	ValidateUser(user *model.User) error
}

// DefaultUserValidator implements Validator for user data validation
type DefaultUserValidator struct{}

// NewDefaultUserValidator creates a new instance of DefaultUserValidator
func NewDefaultUserValidator() *DefaultUserValidator {
	return &DefaultUserValidator{}
}

// ValidateUser validates user fields
func (v *DefaultUserValidator) ValidateUser(user *model.User) error {
	if user.FirstName == "" {
		return errors.New("first_name is required")
	}
	if user.LastName == "" {
		return errors.New("last_name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New("invalid email format")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if user.Country == "" {
		return errors.New("country is required")
	}
	return nil
}
