package security

import "golang.org/x/crypto/bcrypt"

// BcryptPasswordHasher implements PasswordHasher using bcrypt
type BcryptPasswordHasher struct{}

// NewBcryptPasswordHasher creates a new instance of BcryptPasswordHasher
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

// HashPassword hashes a password using bcrypt
func (b *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword checks if a hashed password matches the plain password
func (b *BcryptPasswordHasher) ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
