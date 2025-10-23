package orm

type DB struct {
	executor *executorWrapper
}

// Open initializes and returns a new DB instance
func Open(driver, dsn string) (*DB, error) {
	execWrapper, err := newExecutorWrapper(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{executor: execWrapper}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	if db.executor == nil || db.executor.exec == nil {
		return nil
	}
	return db.executor.exec.Close()
}

// Table initializes a new query for the specified table
func (db *DB) Table(name string) *Query {
	return &Query{
		table:    name,
		executor: db.executor,
	}
}
