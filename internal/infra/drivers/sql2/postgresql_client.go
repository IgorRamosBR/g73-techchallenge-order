package sql2

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type sqlClient struct {
	db *sqlx.DB
}

func NewPostgresSQLClient(username, password, host, port, dbname string) (SQLClient, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s host=%s port=%s", strings.TrimSpace(username), dbname, strings.TrimSpace(password), host, port)
	db, err := sqlx.Open("postgres", connStr)
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

func (client sqlClient) Ping() error {
	return client.db.Ping()
}
