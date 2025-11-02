package postgres

import (
	"context"
	"fmt"

	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
)

// internalClient implements Client interface.
type internalClient struct {
	db     *pg.DB
	logger logger.Logger
}

type transactionContextKeyType string

var transactionContextKey = transactionContextKeyType("transaction")

type QueryResult struct {
	RowsAffected int
	RowsReturned int
}

// Client will manage the infrastructure to deal with some database technology.
type Client interface {
	Connect(connectionString string) error
	Close()
	ExecuteCommand(context context.Context, query string, params ...interface{}) (*QueryResult, error)
	ExecuteQuery(context context.Context, output interface{}, query string, params ...interface{}) (*QueryResult, error)
	StartContextualTransaction(context context.Context) (Transaction, context.Context, error)
}

func NewClient(logger logger.Logger) Client {
	return &internalClient{db: nil, logger: logger}
}

func (client *internalClient) Connect(connectionString string) error {
	options, err := pg.ParseURL(connectionString)

	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	client.db = pg.Connect(options)

	err = client.db.Ping(context.Background())

	if err != nil {
		return fmt.Errorf("unable to ping database and check connectivity: %w", err)
	}

	return nil
}

func (client *internalClient) Close() {
	client.db.Close()
}

func (client *internalClient) StartContextualTransaction(ctx context.Context) (Transaction, context.Context, error) {
	var txWrapper *internalTransaction
	txWrapperAsInterface := ctx.Value(transactionContextKey)

	if txWrapperAsInterface != nil {
		txWrapper, ok := txWrapperAsInterface.(*internalTransaction)
		if !ok {
			return nil, nil, fmt.Errorf("failed to start transaction with context: unable to get transaction")
		}

		return txWrapper, ctx, nil
	}

	tx, err := client.db.BeginContext(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction with context: %w", err)
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction with context: %w", err)
	}

	txWrapper = &internalTransaction{tx: tx, logger: client.logger}
	txContext := context.WithValue(ctx, transactionContextKey, txWrapper)

	return txWrapper, txContext, nil
}

func (client *internalClient) ExecuteQuery(context context.Context, output interface{}, query string, params ...interface{}) (*QueryResult, error) {
	txWrapperAsInterface := context.Value(transactionContextKey)
	var err error
	var result orm.Result

	wrappedParams := client.wrapParamsIfNeeded(params)

	if txWrapperAsInterface != nil {
		txWrapper, ok := txWrapperAsInterface.(*internalTransaction)
		if !ok {
			return nil, fmt.Errorf("failed to execute query '%s' with context: unable to get transaction", query)
		}

		result, err = txWrapper.tx.QueryContext(context, output, query, wrappedParams...)
	} else {
		result, err = client.db.QueryContext(context, output, query, wrappedParams...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query '%s' with context: %w", query, err)
	}

	return &QueryResult{
		RowsAffected: result.RowsAffected(),
		RowsReturned: result.RowsReturned(),
	}, nil
}

func (client *internalClient) ExecuteCommand(context context.Context, query string, params ...interface{}) (*QueryResult, error) {
	txWrapperAsInterface := context.Value(transactionContextKey)
	var err error
	var result orm.Result

	wrappedParams := client.wrapParamsIfNeeded(params)

	if txWrapperAsInterface != nil {
		txWrapper, ok := txWrapperAsInterface.(*internalTransaction)
		if !ok {
			return nil, fmt.Errorf("failed to execute command '%s' with context: unable to get transaction", query)
		}

		result, err = txWrapper.tx.ExecContext(context, query, wrappedParams...)
	} else {
		result, err = client.db.ExecContext(context, query, wrappedParams...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query '%s' with context: %w", query, err)
	}

	return &QueryResult{
		RowsAffected: result.RowsAffected(),
		RowsReturned: result.RowsReturned(),
	}, nil
}

func (client *internalClient) wrapParamsIfNeeded(params []interface{}) []interface{} {
	for i := 0; i < len(params); i++ {
		switch params[i].(type) {
		case []int:
			params[i] = pg.In(params[i])
		case []string:
			params[i] = pg.In(params[i])
		}
	}

	return params
}
