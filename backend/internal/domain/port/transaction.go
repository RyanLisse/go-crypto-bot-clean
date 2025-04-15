package port

import (
	"context"
)

// TxContextKey is the key used to store the transaction in the context
type txContextKey struct{}

// TxContextKey is the key used to store the transaction in the context
var TxContextKey = txContextKey{}

// TransactionManager defines the interface for transaction management
type TransactionManager interface {
	// WithTransaction executes the given function within a transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
