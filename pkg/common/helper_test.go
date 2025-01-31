package common

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetCorrelationID_ValidString ensures it retrieves a string Correlation ID correctly
func TestGetCorrelationID_ValidString(t *testing.T) {
	expectedID := "12345-abcde"
	ctx := context.WithValue(context.Background(), CorrelationID, expectedID)

	result := GetCorrelationID(ctx)

	assert.Equal(t, expectedID, result, "Expected correlation ID to match input string")
}

// TestGetCorrelationID_ValidUUID ensures it retrieves a UUID Correlation ID correctly
func TestGetCorrelationID_ValidUUID(t *testing.T) {
	expectedUUID := uuid.New()
	ctx := context.WithValue(context.Background(), CorrelationID, expectedUUID)

	result := GetCorrelationID(ctx)

	assert.Equal(t, expectedUUID.String(), result, "Expected correlation ID to match input UUID")
}

// TestGetCorrelationID_InvalidType ensures it returns an empty string for incorrect types
func TestGetCorrelationID_InvalidType(t *testing.T) {
	ctx := context.WithValue(context.Background(), CorrelationID, 12345) // Integer instead of string/UUID

	result := GetCorrelationID(ctx)

	assert.Equal(t, "", result, "Expected empty string for invalid correlation ID type")
}

// TestGetCorrelationID_NilContext ensures it handles nil context correctly
func TestGetCorrelationID_NilContext(t *testing.T) {
	result := GetCorrelationID(nil)

	assert.Equal(t, "", result, "Expected empty string when context is nil")
}
