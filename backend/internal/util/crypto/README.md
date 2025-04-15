# Encryption Utilities

This package provides encryption utilities for securely handling sensitive data in the wallet authentication system.

## Overview

The encryption utilities provide a secure way to handle sensitive data such as API credentials, configuration values, and environment variables. The utilities use industry-standard encryption algorithms and best practices to ensure the security of sensitive data.

## Components

### 1. Encryption Service

The encryption service provides a simple interface for encrypting and decrypting data. It supports two types of encryption services:

- **Basic Encryption Service**: A simple encryption service that uses AES-256-GCM for encryption and decryption.
- **Enhanced Encryption Service**: An advanced encryption service that supports key rotation and additional data types.

### 2. Key Manager

The key manager handles encryption keys, including key rotation and secure storage. It supports loading keys from environment variables and provides a simple interface for key management.

### 3. Config Manager

The config manager provides a secure way to store and retrieve configuration values. It encrypts the entire configuration file and provides a simple interface for managing configuration values.

### 4. Environment Variable Manager

The environment variable manager provides a secure way to store and retrieve environment variables. It supports encrypting and decrypting environment variable files and loading encrypted environment variables.

## Usage

### Encryption Service

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

### Key Manager

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

### Config Manager

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

### Environment Variable Manager

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
```

## Security Considerations

- Encryption keys should be stored securely and never committed to version control.
- Use environment variables or a secure key management service to store encryption keys.
- Rotate encryption keys regularly to minimize the impact of key compromise.
- The encryption utilities use AES-256-GCM, which is a secure encryption algorithm.
- The nonce (initialization vector) is generated randomly for each encryption operation.
- The authentication tag is included in the ciphertext to ensure data integrity.

## Best Practices

1. **Use the Enhanced Encryption Service**: The enhanced encryption service provides additional security features such as key rotation and support for different data types.

2. **Rotate Keys Regularly**: Regularly rotate encryption keys to minimize the impact of key compromise.

3. **Encrypt Sensitive Data**: Always encrypt sensitive data before storage or transmission.

4. **Use Secure Key Storage**: Store encryption keys securely and never commit them to version control.

5. **Validate Input**: Validate input data before encryption to prevent security vulnerabilities.

## For More Information

See the [Encryption Utilities Documentation](../../../docs/encryption.md) for more detailed information.
