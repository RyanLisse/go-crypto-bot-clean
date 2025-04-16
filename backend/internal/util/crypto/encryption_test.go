package crypto

import (
	"encoding/base64"
	"os"
	"testing"
	"time"
)

func TestBasicEncryption(t *testing.T) {
	// Set up test key
	testKey := "Wn3PvhLOYk0QpFdod9qUDRRik9cI8jD3noi0TgrTJ1M="
	os.Setenv("ENCRYPTION_KEY", testKey)
	defer os.Unsetenv("ENCRYPTION_KEY")

	// Create encryption service
	factory, err := NewEncryptionServiceFactory()
	if err != nil {
		t.Fatalf("Failed to create encryption service factory: %v", err)
	}

	service, err := factory.GetEncryptionService(BasicEncryptionService)
	if err != nil {
		t.Fatalf("Failed to get basic encryption service: %v", err)
	}

	// Test encryption and decryption
	plaintext := "This is a test message"
	ciphertext, err := service.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := service.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original: got %q, want %q", decrypted, plaintext)
	}
}

func TestEnhancedEncryption(t *testing.T) {
	// Set up test keys
	keyID := "test-key-1"
	testKey := "Wn3PvhLOYk0QpFdod9qUDRRik9cI8jD3noi0TgrTJ1M="
	os.Setenv("ENCRYPTION_CURRENT_KEY_ID", keyID)
	os.Setenv("ENCRYPTION_KEYS", keyID+":"+testKey)
	defer func() {
		os.Unsetenv("ENCRYPTION_CURRENT_KEY_ID")
		os.Unsetenv("ENCRYPTION_KEYS")
	}()

	// Create encryption service
	factory, err := NewEncryptionServiceFactory()
	if err != nil {
		t.Fatalf("Failed to create encryption service factory: %v", err)
	}

	service, err := factory.GetEncryptionService(EnhancedEncryptionServiceType)
	if err != nil {
		t.Fatalf("Failed to get enhanced encryption service: %v", err)
	}

	// Test encryption and decryption
	plaintext := "This is a test message for enhanced encryption"
	ciphertext, err := service.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := service.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original: got %q, want %q", decrypted, plaintext)
	}
}

func TestKeyRotation(t *testing.T) {
	// Set up initial test keys
	keyID1 := "test-key-1"
	testKey1 := "Wn3PvhLOYk0QpFdod9qUDRRik9cI8jD3noi0TgrTJ1M="
	os.Setenv("ENCRYPTION_CURRENT_KEY_ID", keyID1)
	os.Setenv("ENCRYPTION_KEYS", keyID1+":"+testKey1)
	defer func() {
		os.Unsetenv("ENCRYPTION_CURRENT_KEY_ID")
		os.Unsetenv("ENCRYPTION_KEYS")
	}()

	// Create encryption service
	factory, err := NewEncryptionServiceFactory()
	if err != nil {
		t.Fatalf("Failed to create encryption service factory: %v", err)
	}

	service, err := factory.GetEncryptionService(EnhancedEncryptionServiceType)
	if err != nil {
		t.Fatalf("Failed to get enhanced encryption service: %v", err)
	}

	// Encrypt with first key
	plaintext := "This is a test message for key rotation"
	ciphertext1, err := service.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt with first key: %v", err)
	}

	// Rotate key
	keyID2 := "test-key-2"
	testKey2 := "XmFP8d2KzHIhDFGW0rAqPzlJ3QTYbN5UvCxE6sR4o7w="
	os.Setenv("ENCRYPTION_CURRENT_KEY_ID", keyID2)
	os.Setenv("ENCRYPTION_KEYS", keyID2+":"+testKey2+","+keyID1+":"+testKey1)

	// Create new factory and service with rotated keys
	factory2, err := NewEncryptionServiceFactory()
	if err != nil {
		t.Fatalf("Failed to create encryption service factory after key rotation: %v", err)
	}

	service2, err := factory2.GetEncryptionService(EnhancedEncryptionServiceType)
	if err != nil {
		t.Fatalf("Failed to get enhanced encryption service after key rotation: %v", err)
	}

	// Encrypt with second key
	ciphertext2, err := service2.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt with second key: %v", err)
	}

	// Decrypt both ciphertexts with new service
	decrypted1, err := service2.Decrypt(ciphertext1)
	if err != nil {
		t.Fatalf("Failed to decrypt first ciphertext after key rotation: %v", err)
	}

	decrypted2, err := service2.Decrypt(ciphertext2)
	if err != nil {
		t.Fatalf("Failed to decrypt second ciphertext: %v", err)
	}

	if decrypted1 != plaintext {
		t.Errorf("Decrypted text from first key does not match original: got %q, want %q", decrypted1, plaintext)
	}

	if decrypted2 != plaintext {
		t.Errorf("Decrypted text from second key does not match original: got %q, want %q", decrypted2, plaintext)
	}
}

func TestEncryptionTypes(t *testing.T) {
	// Set up test keys
	keyID := "test-key-1"
	testKey := "Wn3PvhLOYk0QpFdod9qUDRRik9cI8jD3noi0TgrTJ1M="
	os.Setenv("ENCRYPTION_CURRENT_KEY_ID", keyID)
	os.Setenv("ENCRYPTION_KEYS", keyID+":"+testKey)
	os.Setenv("ENCRYPTION_KEY", testKey)
	defer func() {
		os.Unsetenv("ENCRYPTION_CURRENT_KEY_ID")
		os.Unsetenv("ENCRYPTION_KEYS")
		os.Unsetenv("ENCRYPTION_KEY")
	}()

	// Create encryption service
	factory, err := NewEncryptionServiceFactory()
	if err != nil {
		t.Fatalf("Failed to create encryption service factory: %v", err)
	}

	service, err := factory.GetEncryptionService(EnhancedEncryptionServiceType)
	if err != nil {
		t.Fatalf("Failed to get enhanced encryption service: %v", err)
	}

	enhancedService, ok := service.(*EnhancedEncryptionService)
	if !ok {
		t.Fatalf("Service is not an EnhancedEncryptionService")
	}

	// Test JSON encryption
	type TestStruct struct {
		Name    string    `json:"name"`
		Age     int       `json:"age"`
		Created time.Time `json:"created"`
	}

	testData := TestStruct{
		Name:    "Test User",
		Age:     30,
		Created: time.Now(),
	}

	ciphertext, err := enhancedService.EncryptJSON(testData)
	if err != nil {
		t.Fatalf("Failed to encrypt JSON: %v", err)
	}

	var decrypted TestStruct
	err = enhancedService.DecryptJSON(ciphertext, &decrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt JSON: %v", err)
	}

	if decrypted.Name != testData.Name || decrypted.Age != testData.Age {
		t.Errorf("Decrypted JSON does not match original: got %+v, want %+v", decrypted, testData)
	}

	// Test integer encryption
	intValue := int64(12345)
	intCiphertext, err := enhancedService.EncryptInt(intValue)
	if err != nil {
		t.Fatalf("Failed to encrypt integer: %v", err)
	}

	decryptedInt, err := enhancedService.DecryptInt(intCiphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt integer: %v", err)
	}

	if decryptedInt != intValue {
		t.Errorf("Decrypted integer does not match original: got %d, want %d", decryptedInt, intValue)
	}

	// Test float encryption
	floatValue := 123.45
	floatCiphertext, err := enhancedService.EncryptFloat(floatValue)
	if err != nil {
		t.Fatalf("Failed to encrypt float: %v", err)
	}

	decryptedFloat, err := enhancedService.DecryptFloat(floatCiphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt float: %v", err)
	}

	if decryptedFloat != floatValue {
		t.Errorf("Decrypted float does not match original: got %f, want %f", decryptedFloat, floatValue)
	}
}

func TestKeyGenerator(t *testing.T) {
	generator := NewKeyGenerator()

	// Test key generation
	key, err := generator.GenerateAES256Key()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Decode key
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Fatalf("Failed to decode key: %v", err)
	}

	// Check key length
	if len(keyBytes) != 32 {
		t.Errorf("Key length is incorrect: got %d, want 32", len(keyBytes))
	}

	// Test key pair generation
	keyID, key, err := generator.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if keyID == "" {
		t.Errorf("Key ID is empty")
	}

	// Test key config generation
	config, err := generator.GenerateKeyConfig()
	if err != nil {
		t.Fatalf("Failed to generate key config: %v", err)
	}

	if config["ENCRYPTION_CURRENT_KEY_ID"] == "" {
		t.Errorf("ENCRYPTION_CURRENT_KEY_ID is empty")
	}

	if config["ENCRYPTION_KEYS"] == "" {
		t.Errorf("ENCRYPTION_KEYS is empty")
	}
}
