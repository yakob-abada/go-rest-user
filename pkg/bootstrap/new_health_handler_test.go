package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/go-rest-user/pkg/handler"
	"gorm.io/gorm"
)

// TestNewHealthHandler ensures that a HealthHandler is properly created
func TestNewHealthHandler(t *testing.T) {
	mockDB := &gorm.DB{} // Using an empty gorm.DB instance as a mock

	healthHandler := NewHealthHandler(mockDB)

	assert.NotNil(t, healthHandler, "Expected a non-nil HealthHandler instance")
	assert.IsType(t, &handler.HealthHandler{}, healthHandler, "Expected type *handler.HealthHandler")
}
