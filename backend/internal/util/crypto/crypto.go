package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var (
	encryptionKey []byte
)

func init() {
	keyB64 := os.Getenv("MEXC_CRED_ENCRYPTION_KEY")
	if keyB64 == "" {
		// Always use default key for development and testing
		keyB64 = "Wn3PvhLOYk0QpFdod9qUDRRik9cI8jD3noi0TgrTJ1M="
	}

	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil || len(key) != 32 {
		// Fail explicitly rather than using an insecure zero key
		panic("MEXC_CRED_ENCRYPTION_KEY must be a valid 32-byte key encoded as base64")
	}
	encryptionKey = key
}

// Encrypt encrypts a string using AES-GCM
func Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	// Never use more than 2^32 random nonces with a given key
	nonce := make([]byte, 12)
	if _, readErr := io.ReadFull(rand.Reader, nonce); readErr != nil {
		return "", readErr
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Prepend nonce to ciphertext
	result := make([]byte, len(nonce)+len(ciphertext))
	copy(result, nonce)
	copy(result[len(nonce):], ciphertext)

	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts a string using AES-GCM
func Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < 12 {
		return "", errors.New("invalid ciphertext")
	}

	nonce := data[:12]
	ciphertext := data[12:]

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
