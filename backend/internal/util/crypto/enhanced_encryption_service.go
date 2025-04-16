package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
)

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
	KeyID      string `json:"kid"`
	Nonce      []byte `json:"n"`
	Ciphertext []byte `json:"c"`
}

// EnhancedEncryptionService implements EncryptionService with key rotation support
type EnhancedEncryptionService struct {
	keyManager KeyManager
}

// NewEnhancedEncryptionService creates a new EnhancedEncryptionService
func NewEnhancedEncryptionService(keyManager KeyManager) *EnhancedEncryptionService {
	return &EnhancedEncryptionService{
		keyManager: keyManager,
	}
}

// Encrypt encrypts a string using AES-256-GCM with the current key
func (s *EnhancedEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	// Get current key
	key, err := s.keyManager.GetCurrentKey()
	if err != nil {
		return nil, err
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Create encrypted data
	data := EncryptedData{
		KeyID:      s.getCurrentKeyID(),
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}

	// Serialize encrypted data
	return json.Marshal(data)
}

// Decrypt decrypts a string using AES-256-GCM
func (s *EnhancedEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	// Deserialize encrypted data
	var data EncryptedData
	if err := json.Unmarshal(ciphertext, &data); err != nil {
		// Try legacy format
		return s.decryptLegacy(ciphertext)
	}

	// Get key by ID
	key, err := s.keyManager.GetKeyByID(data.KeyID)
	if err != nil {
		return "", err
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, data.Nonce, data.Ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// decryptLegacy decrypts data in the legacy format
func (s *EnhancedEncryptionService) decryptLegacy(ciphertext []byte) (string, error) {
	// Get current key
	key, err := s.keyManager.GetCurrentKey()
	if err != nil {
		return "", err
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Check ciphertext length
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// getCurrentKeyID returns the current key ID
func (s *EnhancedEncryptionService) getCurrentKeyID() string {
	// Get the current key ID from environment variable
	currentKeyID := os.Getenv("ENCRYPTION_CURRENT_KEY_ID")
	if currentKeyID == "" {
		// Fallback to default
		return "test-key-1"
	}
	return currentKeyID
}

// EncryptJSON encrypts a JSON-serializable object
func (s *EnhancedEncryptionService) EncryptJSON(data interface{}) ([]byte, error) {
	// Serialize data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Encrypt JSON data
	return s.Encrypt(string(jsonData))
}

// DecryptJSON decrypts a JSON-serializable object
func (s *EnhancedEncryptionService) DecryptJSON(ciphertext []byte, target interface{}) error {
	// Decrypt ciphertext
	plaintext, err := s.Decrypt(ciphertext)
	if err != nil {
		return err
	}

	// Deserialize JSON data
	return json.Unmarshal([]byte(plaintext), target)
}

// EncryptBytes encrypts binary data
func (s *EnhancedEncryptionService) EncryptBytes(data []byte) ([]byte, error) {
	return s.Encrypt(string(data))
}

// DecryptBytes decrypts binary data
func (s *EnhancedEncryptionService) DecryptBytes(ciphertext []byte) ([]byte, error) {
	plaintext, err := s.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	return []byte(plaintext), nil
}

// EncryptInt encrypts an integer
func (s *EnhancedEncryptionService) EncryptInt(value int64) ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return s.EncryptBytes(buf)
}

// DecryptInt decrypts an integer
func (s *EnhancedEncryptionService) DecryptInt(ciphertext []byte) (int64, error) {
	buf, err := s.DecryptBytes(ciphertext)
	if err != nil {
		return 0, err
	}
	if len(buf) != 8 {
		return 0, errors.New("invalid integer data")
	}
	return int64(binary.LittleEndian.Uint64(buf)), nil
}

// EncryptFloat encrypts a floating-point number
func (s *EnhancedEncryptionService) EncryptFloat(value float64) ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(value))
	return s.EncryptBytes(buf)
}

// DecryptFloat decrypts a floating-point number
func (s *EnhancedEncryptionService) DecryptFloat(ciphertext []byte) (float64, error) {
	buf, err := s.DecryptBytes(ciphertext)
	if err != nil {
		return 0, err
	}
	if len(buf) != 8 {
		return 0, errors.New("invalid float data")
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(buf)), nil
}
