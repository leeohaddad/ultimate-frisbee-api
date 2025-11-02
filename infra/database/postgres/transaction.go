package postgres

import (
	"context"
	"fmt"

	pg "github.com/go-pg/pg/v10"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
)

// internalTransaction implements Transaction interface.
type internalTransaction struct {
	tx     *pg.Tx
	logger logger.Logger
}

// Transaction defines an interface to deal with a transaction running into the database.
type Transaction interface {
	Commit(context context.Context) error
	Rollback(context context.Context) error
	Close(context context.Context)
}

func (transaction *internalTransaction) Commit(context context.Context) error {
	err := transaction.tx.CommitContext(context)

	if err != nil {
		return fmt.Errorf("failed to commit transaction with context: %w", err)
	}

	return nil
}

func (transaction *internalTransaction) Rollback(context context.Context) error {
	err := transaction.tx.RollbackContext(context)

	if err != nil {
		return fmt.Errorf("failed to rollback transaction with context: %w", err)
	}

	return nil
}

func (transaction *internalTransaction) Close(context context.Context) {
	err := transaction.tx.CloseContext(context)

	if err != nil {
		transaction.logger.Error(fmt.Sprintf("failed to close transaction with context: %s", err))
	}
}
