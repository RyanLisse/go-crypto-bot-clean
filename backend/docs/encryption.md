# Encryption Utilities Documentation

This document provides comprehensive documentation for the encryption utilities used in the wallet authentication system.

## Overview

The encryption utilities provide a secure way to handle sensitive data such as API credentials, configuration values, and environment variables. The utilities use industry-standard encryption algorithms and best practices to ensure the security of sensitive data.

## Key Components

### 1. Encryption Service

The encryption service provides a simple interface for encrypting and decrypting data. It supports two types of encryption services:

- **Basic Encryption Service**: A simple encryption service that uses AES-256-GCM for encryption and decryption.
- **Enhanced Encryption Service**: An advanced encryption service that supports key rotation and additional data types.

#### Usage

```go
// Create encryption service factory
factory, err := crypto.NewEncryptionServiceFactory()
if err != nil {
    // Handle error
}

// Get basic encryption service
basicSvc, err := factory.GetEncryptionService(crypto.BasicEncryptionService)
if err != nil {
    // Handle error
}

// Get enhanced encryption service
enhancedSvc, err := factory.GetEncryptionService(crypto.EnhancedEncryptionServiceType)
if err != nil {
    // Handle error
}

// Encrypt data
ciphertext, err := basicSvc.Encrypt("sensitive data")
if err != nil {
    // Handle error
}

// Decrypt data
plaintext, err := basicSvc.Decrypt(ciphertext)
if err != nil {
    // Handle error
}
```

### 2. Key Manager

The key manager handles encryption keys, including key rotation and secure storage. It supports loading keys from environment variables and provides a simple interface for key management.

#### Usage

```go
// Create key manager
keyManager, err := crypto.NewEnvKeyManager()
if err != nil {
    // Handle error
}

// Get current key
key, err := keyManager.GetCurrentKey()
if err != nil {
    // Handle error
}

// Rotate key
newKeyID, err := keyManager.RotateKey()
if err != nil {
    // Handle error
}
```

### 3. Credential Manager

The credential manager provides a secure way to store and retrieve API credentials. It encrypts sensitive data before storage and decrypts it when needed.

#### Usage

```go
// Create credential manager
credentialManager := service.NewCredentialManager(credentialRepo, encryptionSvc)

// Store credential
credential, err := credentialManager.StoreCredential(ctx, userID, exchange, apiKey, apiSecret, label)
if err != nil {
    // Handle error
}

// Get credential
credential, err := credentialManager.GetCredential(ctx, id)
if err != nil {
    // Handle error
}

// Get credential with secret
credential, err := credentialManager.GetCredentialWithSecret(ctx, id)
if err != nil {
    // Handle error
}
```

### 4. Config Manager

The config manager provides a secure way to store and retrieve configuration values. It encrypts the entire configuration file and provides a simple interface for managing configuration values.

#### Usage

```go
// Create config manager
configManager, err := crypto.NewConfigManager(encryptionSvc, configPath)
if err != nil {
    // Handle error
}

// Set value
err := configManager.SetValue("key", "value")
if err != nil {
    // Handle error
}

// Get value
value, err := configManager.GetValue("key")
if err != nil {
    // Handle error
}
```

### 5. Environment Variable Manager

The environment variable manager provides a secure way to store and retrieve environment variables. It supports encrypting and decrypting environment variable files and loading encrypted environment variables.

#### Usage

```go
// Create environment variable manager
envManager := crypto.NewEnvManager(encryptionSvc, envPath)

// Load environment variables
err := envManager.LoadEnv()
if err != nil {
    // Handle error
}

// Save environment variables
err := envManager.SaveEnv(vars, true)
if err != nil {
    // Handle error
}

// Encrypt environment variable file
err := envManager.EncryptEnvFile(inputPath, outputPath)
if err != nil {
    // Handle error
}

// Decrypt environment variable file
err := envManager.DecryptEnvFile(inputPath, outputPath)
if err != nil {
    // Handle error
}
```

## Command-Line Tools

### 1. Key Generator

The key generator tool provides a command-line interface for generating encryption keys.

#### Usage

```bash
# Generate a new encryption key
go run cmd/keygen/main.go -generate

# Generate a new encryption key with specific bits
go run cmd/keygen/main.go -generate -bits 256

# Generate a new encryption key configuration
go run cmd/keygen/main.go -generate -env

# Rotate encryption keys
go run cmd/keygen/main.go -rotate
```

### 2. Environment Variable Tool

The environment variable tool provides a command-line interface for encrypting and decrypting environment variable files.

#### Usage

```bash
# Encrypt an environment variable file
go run cmd/envtool/main.go -encrypt -input .env -output .env.enc -key <encryption-key>

# Decrypt an environment variable file
go run cmd/envtool/main.go -decrypt -input .env.enc -output .env -key <encryption-key>
```

## Security Considerations

### 1. Key Management

- Encryption keys should be stored securely and never committed to version control.
- Use environment variables or a secure key management service to store encryption keys.
- Rotate encryption keys regularly to minimize the impact of key compromise.

### 2. Encryption Algorithm

- The encryption utilities use AES-256-GCM, which is a secure encryption algorithm.
- The nonce (initialization vector) is generated randomly for each encryption operation.
- The authentication tag is included in the ciphertext to ensure data integrity.

### 3. Key Rotation

- The enhanced encryption service supports key rotation, which allows for changing encryption keys without losing access to previously encrypted data.
- When a key is rotated, new data is encrypted with the new key, but old data can still be decrypted with the old key.

### 4. Secure Storage

- Sensitive data should be encrypted before storage.
- Encrypted data should be stored in a secure location.
- Access to encrypted data should be restricted to authorized users.

## Best Practices

1. **Use the Enhanced Encryption Service**: The enhanced encryption service provides additional security features such as key rotation and support for different data types.

2. **Rotate Keys Regularly**: Regularly rotate encryption keys to minimize the impact of key compromise.

3. **Encrypt Sensitive Data**: Always encrypt sensitive data before storage or transmission.

4. **Use Secure Key Storage**: Store encryption keys securely and never commit them to version control.

5. **Validate Input**: Validate input data before encryption to prevent security vulnerabilities.

6. **Handle Errors Properly**: Handle encryption and decryption errors properly to prevent security vulnerabilities.

7. **Use Secure Random Number Generation**: Use secure random number generation for cryptographic operations.

8. **Keep Dependencies Updated**: Keep encryption libraries and dependencies updated to ensure security.

## Troubleshooting

### 1. Encryption Key Issues

- **Invalid Key Length**: Ensure that the encryption key is 32 bytes (256 bits) for AES-256-GCM.
- **Key Not Found**: Ensure that the encryption key is available in the environment or configuration.

### 2. Encryption/Decryption Issues

- **Invalid Ciphertext**: Ensure that the ciphertext is valid and has not been tampered with.
- **Decryption Failed**: Ensure that the correct key is used for decryption.

### 3. Environment Variable Issues

- **Environment Variable Not Found**: Ensure that the environment variable is set before use.
- **Environment Variable File Not Found**: Ensure that the environment variable file exists and is readable.

## References

- [AES-GCM Encryption](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
- [NIST Recommendations for Key Management](https://csrc.nist.gov/publications/detail/sp/800-57-part-1/rev-5/final)
- [Go Cryptography](https://golang.org/pkg/crypto/)
