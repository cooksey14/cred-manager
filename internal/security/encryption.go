package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	models "github.com/cooksey14/cred-manager/internal/models"
)

// Encrypt encrypts plaintext using AES-GCM with the provided key.
func Encrypt(plaintext string, key []byte) (string, string, error) {
	// Ensure the key is 32 bytes for AES-256
	if len(key) != 32 {
		return "", "", fmt.Errorf("encryption key must be 32 bytes long")
	}

	// Create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Generate a random nonce for AES-GCM
	nonce := make([]byte, 12) // AES-GCM recommends a 12-byte nonce
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Create AES-GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Encrypt the plaintext
	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	// Return the ciphertext and nonce as base64-encoded strings
	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(nonce), nil
}

// Decrypt decrypts ciphertext using AES-GCM with the provided key and nonce.
func Decrypt(ciphertext, nonce string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	decodedNonce, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to decode nonce: %v", err)
	}

	plaintext, err := aesGCM.Open(nil, decodedNonce, decodedCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plaintext), nil
}

// Function to save the vault
func saveVault(vault models.Vault, key []byte, filepath string) error {
	data, err := json.Marshal(vault)
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %v", err)
	}

	ciphertext, nonce, err := Encrypt(string(data), key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %v", err)
	}

	encryptedData := fmt.Sprintf("%s|%s", ciphertext, nonce)

	return os.WriteFile(filepath, []byte(encryptedData), 0600)
}

func loadVault(key []byte, filepath string) (models.Vault, error) {
	var vault models.Vault

	encryptedData, err := os.ReadFile(filepath)
	if err != nil {
		return vault, fmt.Errorf("failed to read vault file: %v", err)
	}

	parts := strings.Split(string(encryptedData), "|")
	if len(parts) != 2 {
		return vault, fmt.Errorf("invalid encrypted data format")
	}
	ciphertext := parts[0]
	nonce := parts[1]

	decryptedData, err := Decrypt(ciphertext, nonce, key)
	if err != nil {
		return vault, fmt.Errorf("failed to decrypt vault: %v", err)
	}

	err = json.Unmarshal([]byte(decryptedData), &vault)
	if err != nil {
		return vault, fmt.Errorf("failed to unmarshal vault: %v", err)
	}

	return vault, nil
}
