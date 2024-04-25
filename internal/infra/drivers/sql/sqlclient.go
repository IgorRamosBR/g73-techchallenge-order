package sql

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("entity not found")

type SQLClient interface {
	Find(result any, query string, args ...any) error
	FindOne(result any, query string, args ...any) error
	Exec(query string, args ...any) (ResultWrapper, error)
	ExecWithReturn(query string, args ...any) RowWrapper
	Begin() (TransactionWrapper, error)
	Ping() error
	GetConnection() *sql.DB
}
