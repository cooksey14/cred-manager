package security

import (
	"crypto/rand"
	"math/big"
)

const credentialChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"

// GeneratePassword generates a secure random credential of specified length
func GeneratePassword(length int) (string, error) {
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(credentialChars))))
		if err != nil {
			return "", err
		}
		password[i] = credentialChars[num.Int64()]
	}
	return string(password), nil
}
