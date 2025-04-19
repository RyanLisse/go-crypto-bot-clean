package port

// EncryptionService defines the interface for encryption operations
type EncryptionService interface {
	// Encrypt encrypts a plaintext string
	Encrypt(plaintext string) (string, error)
	
	// Decrypt decrypts a ciphertext string
	Decrypt(ciphertext string) (string, error)
}

// EnhancedEncryptionService extends EncryptionService with key rotation support
type EnhancedEncryptionService interface {
	EncryptionService
	
	// ReEncryptWithCurrentKey re-encrypts data with the current key
	ReEncryptWithCurrentKey(ciphertext string) (string, error)
}
