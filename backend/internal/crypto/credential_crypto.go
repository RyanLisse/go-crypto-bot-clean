package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptCredentialSecret encrypts the secret using the active key and returns ciphertext + key version
func EncryptCredentialSecret(plain string, registry *KeyRegistry) (ciphertext string, keyVersion string, err error) {
	meta, err := registry.GetActiveKey()
	if err != nil {
		return "", "", fmt.Errorf("no active key: %w", err)
	}
	key := meta.Key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", fmt.Errorf("aes.NewCipher: %w", err)
	}
	b := []byte(plain)
	ciphertextBytes := make([]byte, aes.BlockSize+len(b))
	iv := ciphertextBytes[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", fmt.Errorf("iv: %w", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes[aes.BlockSize:], b)
	return base64.StdEncoding.EncodeToString(ciphertextBytes), meta.Version, nil
}

// DecryptCredentialSecret decrypts the secret using the key version from registry
func DecryptCredentialSecret(ciphertext string, keyVersion string, registry *KeyRegistry) (string, error) {
	meta, err := registry.GetKey(keyVersion)
	if err != nil {
		return "", fmt.Errorf("key version lookup: %w", err)
	}
	key := meta.Key
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes.NewCipher: %w", err)
	}
	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertextBytes[:aes.BlockSize]
	b := ciphertextBytes[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(b, b)
	return string(b), nil
}
