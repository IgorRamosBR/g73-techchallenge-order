package sql2

import (
	"errors"
)

var ErrNotFound = errors.New("entity not found")

type SQLClient interface {
	Find(result any, query string, args ...any) error
	FindOne(result any, query string, args ...any) error
	Ping() error
}
