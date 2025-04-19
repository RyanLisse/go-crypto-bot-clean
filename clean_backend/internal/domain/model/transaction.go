package model

import (
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

// Transaction status
type TransactionStatus string

// Transaction type constants
const (
	TransactionTypeDeposit      TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal   TransactionType = "WITHDRAWAL"
	TransactionTypeTrade        TransactionType = "TRADE"
	TransactionTypeFee          TransactionType = "FEE"
	TransactionTypeTransfer     TransactionType = "TRANSFER"
	TransactionTypeDistribution TransactionType = "DISTRIBUTION"
)

// Transaction status constants
const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID        uint64            
	UserID    string            
	Type      TransactionType   
	Asset     Asset             
	Amount    float64           
	Fee       float64           
	Timestamp time.Time         
	Status    TransactionStatus 
	TxID      string            
}
