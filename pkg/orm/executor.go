package orm

import (
	"github.com/i-sub135/i-sub-orm/internal/executor"
)

// executorWrapper is a wrapper around the executor.Executor struct
type executorWrapper struct {
	exec   *executor.Executor
	driver string
}

// newExecutorWrapper creates a new executorWrapper instance
func newExecutorWrapper(driver, dsn string) (*executorWrapper, error) {
	exec, err := executor.NewExecutor(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &executorWrapper{
		exec:   exec,
		driver: driver,
	}, nil
}
