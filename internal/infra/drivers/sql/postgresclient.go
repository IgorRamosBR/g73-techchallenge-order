package sql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type sqlClient struct {
	db *sqlx.DB
}

func NewPostgresSQLClient(username, password, host, port, dbname, sslmode string) (SQLClient, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", strings.TrimSpace(username), strings.TrimSpace(password), host, port, dbname, sslmode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return sqlClient{
		db,
	}, nil
}

func (client sqlClient) Find(result any, query string, args ...any) error {
	return client.db.Select(result, query, args...)
}

func (client sqlClient) FindOne(result any, query string, args ...any) error {
	return client.db.Get(result, query, args...)
}

func (client sqlClient) Exec(query string, args ...any) (ResultWrapper, error) {
	result, err := client.db.Exec(query, args...)
	return NewResultWrapper(result), err
}

func (client sqlClient) ExecWithReturn(query string, args ...any) RowWrapper {
	return NewRowWrapper(client.db.QueryRow(query, args...))
}

func (client sqlClient) Begin() (TransactionWrapper, error) {
	tx, err := client.db.Begin()
	return NewTransactionWrapper(tx), err
}

func (client sqlClient) Ping() error {
	err := client.db.Ping()
	return err
}

func (client sqlClient) GetConnection() *sql.DB {
	return client.db.DB
}
