package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/go-rest-user/pkg/handler"
	"github.com/yakob-abada/go-rest-user/pkg/logging"
	"github.com/yakob-abada/go-rest-user/pkg/publisher"
	"gorm.io/gorm"
)

// TestNewUserHandler ensures that NewUserHandler correctly initializes a UserHandler instance
func TestNewUserHandler(t *testing.T) {
	mockDB := &gorm.DB{}
	mockPublisher := new(publisher.MockPublisher)
	mockLogger := new(logging.MockLogger)

	userHandler := NewUserHandler(mockDB, mockPublisher, mockLogger)

	// Assertions
	assert.NotNil(t, userHandler, "Expected a non-nil UserHandler instance")
	assert.IsType(t, &handler.UserHandler{}, userHandler, "Expected type *handler.UserHandler")
}
