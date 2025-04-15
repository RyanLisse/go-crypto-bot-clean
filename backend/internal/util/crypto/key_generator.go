package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// KeyGenerator generates secure encryption keys
type KeyGenerator struct{}

// NewKeyGenerator creates a new KeyGenerator
func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{}
}

// GenerateKey generates a new random encryption key
func (g *KeyGenerator) GenerateKey(bits int) (string, error) {
	if bits%8 != 0 {
		return "", fmt.Errorf("bits must be a multiple of 8")
	}

	bytes := bits / 8
	key := make([]byte, bytes)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateAES256Key generates a new random AES-256 encryption key
func (g *KeyGenerator) GenerateAES256Key() (string, error) {
	return g.GenerateKey(256)
}

// GenerateKeyPair generates a new key ID and key
func (g *KeyGenerator) GenerateKeyPair() (string, string, error) {
	// Generate key ID
	keyID := generateKeyID()

	// Generate key
	key, err := g.GenerateAES256Key()
	if err != nil {
		return "", "", err
	}

	return keyID, key, nil
}

// GenerateKeyConfig generates a complete key configuration for environment variables
func (g *KeyGenerator) GenerateKeyConfig() (map[string]string, error) {
	// Generate key ID and key
	keyID, key, err := g.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	// Create config
	config := make(map[string]string)
	config["ENCRYPTION_CURRENT_KEY_ID"] = keyID
	config["ENCRYPTION_KEYS"] = fmt.Sprintf("%s:%s", keyID, key)

	return config, nil
}

// RotateKeyConfig rotates the keys in an existing configuration
func (g *KeyGenerator) RotateKeyConfig(currentConfig map[string]string) (map[string]string, error) {
	// Get current keys
	currentKeys := currentConfig["ENCRYPTION_KEYS"]
	if currentKeys == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEYS not found in current config")
	}

	// Generate new key ID and key
	keyID, key, err := g.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	// Create new config
	newConfig := make(map[string]string)
	for k, v := range currentConfig {
		newConfig[k] = v
	}

	newConfig["ENCRYPTION_CURRENT_KEY_ID"] = keyID
	newConfig["ENCRYPTION_KEYS"] = fmt.Sprintf("%s:%s,%s", keyID, key, currentKeys)

	return newConfig, nil
}
