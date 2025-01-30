package security

// PasswordHasher defines an interface for password hashing
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, plainPassword string) bool
}
