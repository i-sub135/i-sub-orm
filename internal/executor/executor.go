package executor

import (
	"database/sql"
	"fmt"
)

type Executor struct {
	DB *sql.DB
}

// NewExecutor creates a new Executor with a database connection.
func NewExecutor(driver, dsn string) (*Executor, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection : %w", err)
	}
	return &Executor{DB: db}, nil
}

// Query executes a query and returns the resulting rows.
func (e *Executor) Query(query string, args ...any) (*sql.Rows, error) {
	return e.DB.Query(query, args...)
}

// Close closes the database connection.
func (e *Executor) Close() error {
	return e.DB.Close()
}
