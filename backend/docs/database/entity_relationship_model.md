# Entity-Relationship Model and Data Dictionary

## Overview

This document defines the entity-relationship model and data dictionary for the wallet authentication system. It includes all tables, columns, relationships, and constraints for the database schema.

## Entity-Relationship Diagram

```
+-------------------+       +-------------------+       +-------------------+
|                   |       |                   |       |                   |
|   UserEntity      |       |  APICredential    |       |    WalletEntity   |
|                   |       |                   |       |                   |
+-------------------+       +-------------------+       +-------------------+
| PK: id            |<----->| PK: id            |       | PK: id            |
| email             |       | FK: user_id       |       | FK: user_id       |
| name              |       | exchange          |       | exchange          |
| created_at        |       | api_key           |       | last_updated      |
| updated_at        |       | api_secret        |       | total_usd_value   |
|                   |       | label             |       |                   |
|                   |       | created_at        |       |                   |
|                   |       | updated_at        |       |                   |
+-------------------+       +-------------------+       +-------------------+
                                                                |
                                                                |
                                                                v
                                                        +-------------------+
                                                        |                   |
                                                        |  BalanceEntity    |
                                                        |                   |
                                                        +-------------------+
                                                        | PK: id            |
                                                        | FK: wallet_id     |
                                                        | asset             |
                                                        | free              |
                                                        | locked            |
                                                        | total             |
                                                        | usd_value         |
                                                        | created_at        |
                                                        | updated_at        |
                                                        +-------------------+
```

## Data Dictionary

### UserEntity

Represents a user in the system.

| Column     | Type         | Constraints       | Description                           |
|------------|--------------|-------------------|---------------------------------------|
| id         | VARCHAR(50)  | PK, NOT NULL      | Unique identifier (from Clerk)        |
| email      | VARCHAR(100) | UNIQUE, NOT NULL  | User's email address                  |
| name       | VARCHAR(100) | NULL              | User's full name                      |
| created_at | TIMESTAMP    | NOT NULL          | When the record was created           |
| updated_at | TIMESTAMP    | NOT NULL          | When the record was last updated      |

Indexes:
- PRIMARY KEY (id)
- UNIQUE INDEX idx_user_email (email)

### APICredentialEntity

Stores encrypted API credentials for cryptocurrency exchanges.

| Column     | Type         | Constraints       | Description                           |
|------------|--------------|-------------------|---------------------------------------|
| id         | VARCHAR(50)  | PK, NOT NULL      | Unique identifier (UUID)              |
| user_id    | VARCHAR(50)  | FK, NOT NULL      | Reference to UserEntity.id            |
| exchange   | VARCHAR(20)  | NOT NULL          | Exchange name (e.g., "MEXC")          |
| api_key    | VARCHAR(100) | NOT NULL          | Public API key                        |
| api_secret | BLOB         | NOT NULL          | Encrypted API secret                  |
| label      | VARCHAR(50)  | NULL              | User-defined label                    |
| created_at | TIMESTAMP    | NOT NULL          | When the record was created           |
| updated_at | TIMESTAMP    | NOT NULL          | When the record was last updated      |

Indexes:
- PRIMARY KEY (id)
- INDEX idx_api_credentials_user_id (user_id)
- INDEX idx_api_credentials_exchange (exchange)
- UNIQUE INDEX idx_api_credentials_user_exchange_label (user_id, exchange, label)

Constraints:
- FOREIGN KEY (user_id) REFERENCES UserEntity(id) ON DELETE CASCADE

### WalletEntity

Represents a cryptocurrency wallet with balances.

| Column          | Type         | Constraints       | Description                           |
|-----------------|--------------|-------------------|---------------------------------------|
| id              | INTEGER      | PK, NOT NULL      | Auto-incrementing ID                  |
| user_id         | VARCHAR(50)  | FK, NOT NULL      | Reference to UserEntity.id            |
| exchange        | VARCHAR(20)  | NOT NULL          | Exchange name (e.g., "MEXC")          |
| last_updated    | TIMESTAMP    | NOT NULL          | When wallet data was last updated     |
| total_usd_value | DECIMAL(18,8)| NOT NULL          | Total USD value of all balances       |
| created_at      | TIMESTAMP    | NOT NULL          | When the record was created           |
| updated_at      | TIMESTAMP    | NOT NULL          | When the record was last updated      |

Indexes:
- PRIMARY KEY (id)
- INDEX idx_wallet_user_id (user_id)
- UNIQUE INDEX idx_wallet_user_exchange (user_id, exchange)

Constraints:
- FOREIGN KEY (user_id) REFERENCES UserEntity(id) ON DELETE CASCADE

### BalanceEntity

Represents a balance for a specific asset in a wallet.

| Column     | Type         | Constraints       | Description                           |
|------------|--------------|-------------------|---------------------------------------|
| id         | INTEGER      | PK, NOT NULL      | Auto-incrementing ID                  |
| wallet_id  | INTEGER      | FK, NOT NULL      | Reference to WalletEntity.id          |
| asset      | VARCHAR(20)  | NOT NULL          | Asset symbol (e.g., "BTC")            |
| free       | DECIMAL(18,8)| NOT NULL          | Available balance                     |
| locked     | DECIMAL(18,8)| NOT NULL          | Locked balance                        |
| total      | DECIMAL(18,8)| NOT NULL          | Total balance (free + locked)         |
| usd_value  | DECIMAL(18,8)| NOT NULL          | USD equivalent value                  |
| created_at | TIMESTAMP    | NOT NULL          | When the record was created           |
| updated_at | TIMESTAMP    | NOT NULL          | When the record was last updated      |

Indexes:
- PRIMARY KEY (id)
- INDEX idx_balance_wallet_id (wallet_id)
- INDEX idx_balance_asset (asset)
- UNIQUE INDEX idx_balance_wallet_asset (wallet_id, asset)

Constraints:
- FOREIGN KEY (wallet_id) REFERENCES WalletEntity(id) ON DELETE CASCADE

## Relationships

1. **User to API Credentials**: One-to-Many
   - A user can have multiple API credentials
   - Each API credential belongs to exactly one user

2. **User to Wallets**: One-to-Many
   - A user can have multiple wallets (one per exchange)
   - Each wallet belongs to exactly one user

3. **Wallet to Balances**: One-to-Many
   - A wallet can have multiple balances (one per asset)
   - Each balance belongs to exactly one wallet

## Notes on Security

- API secrets are stored encrypted using AES-256-GCM
- The encryption key is stored in environment variables, not in the database
- User IDs are derived from Clerk authentication, providing an additional layer of security
