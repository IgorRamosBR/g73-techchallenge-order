package sql

import (
	"database/sql"
	"errors"
)

type TransactionWrapper interface {
	FindOne(query string, args ...any) RowWrapper
	Exec(query string, args ...any) (ResultWrapper, error)
	ExecWithReturn(query string, args ...any) RowWrapper
	Commit() error
	Rollback() error
}

type transactionWrapper struct {
	tx *sql.Tx
}

func NewTransactionWrapper(tx *sql.Tx) TransactionWrapper {
	return transactionWrapper{
		tx,
	}
}

func (t transactionWrapper) FindOne(query string, args ...any) RowWrapper {
	row := t.tx.QueryRow(query, args...)
	return NewRowWrapper(row)
}

func (t transactionWrapper) Exec(query string, args ...any) (ResultWrapper, error) {
	result, err := t.tx.Exec(query, args...)
	return result, err
}

func (t transactionWrapper) ExecWithReturn(query string, args ...any) RowWrapper {
	row := t.tx.QueryRow(query, args...)
	return NewRowWrapper(row)
}

func (t transactionWrapper) Commit() error {
	err := t.tx.Commit()
	return err
}

func (t transactionWrapper) Rollback() error {
	err := t.tx.Rollback()
	// skip the log if a transaction has already been committed or rolled back
	if err != nil && errors.Is(err, sql.ErrTxDone) {
		return nil
	}
	return err
}
