//go:build unit_test

package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptPasswordHasher_HashPassword_Success(t *testing.T) {
	hasher := NewBcryptPasswordHasher()
	password := "securepassword123"

	hashedPassword, err := hasher.HashPassword(password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword, "Hashed password should not match the plain password")
	assert.Contains(t, hashedPassword, "$2a$", "Hashed password should be in bcrypt format")
}

func TestBcryptPasswordHasher_HashPassword_DifferentHashes(t *testing.T) {
	hasher := NewBcryptPasswordHasher()
	password := "securepassword123"

	hash1, err1 := hasher.HashPassword(password)
	hash2, err2 := hasher.HashPassword(password)

	// Assertions
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "Hashing the same password twice should produce different hashes")
}

func TestBcryptPasswordHasher_ComparePassword_CorrectPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()
	password := "securepassword123"

	hashedPassword, err := hasher.HashPassword(password)
	assert.NoError(t, err)

	match := hasher.ComparePassword(hashedPassword, password)

	// Assertions
	assert.True(t, match, "Correct password should match hashed password")
}

func TestBcryptPasswordHasher_ComparePassword_WrongPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()
	password := "securepassword123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := hasher.HashPassword(password)
	assert.NoError(t, err)

	match := hasher.ComparePassword(hashedPassword, wrongPassword)

	// Assertions
	assert.False(t, match, "Wrong password should not match hashed password")
}

func TestBcryptPasswordHasher_ComparePassword_InvalidHash(t *testing.T) {
	hasher := NewBcryptPasswordHasher()
	password := "securepassword123"
	invalidHash := "notahash"

	match := hasher.ComparePassword(invalidHash, password)

	// Assertions
	assert.False(t, match, "Invalid hash format should return false")
}
