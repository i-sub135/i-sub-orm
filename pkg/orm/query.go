package orm

import (
	"fmt"
	"strings"

	"github.com/i-sub135/i-sub-orm/internal/expr"
	"github.com/i-sub135/i-sub-orm/internal/utils"
)

type Query struct {
	table    string
	fields   []string
	where    []string
	args     []any
	executor *executorWrapper
}

func (q *Query) Select(cols ...string) *Query {
	q.fields = append(q.fields, cols...)
	return q
}

// Flexible Where(): bisa string atau expr (Eq, Neq, dll)
func (q *Query) Where(cond any, args ...any) *Query {
	switch c := cond.(type) {
	case string:
		q.where = append(q.where, c)
		q.args = append(q.args, args...)
	default:
		sql, a := expr.Compile(c)
		if sql != "" {
			q.where = append(q.where, sql)
			q.args = append(q.args, a...)
		}
	}
	return q
}

func (q *Query) Build() string {
	sql := "SELECT "

	// Handle fields
	if len(q.fields) == 0 {
		sql += "*"
	} else {
		sql += strings.Join(q.fields, ", ")
	}

	// Add FROM clause
	sql += " FROM " + q.table

	// Add WHERE clause
	if len(q.where) > 0 {
		sql += " WHERE " + strings.Join(q.where, " AND ")
	}

	return sql
}

func (q *Query) Get(dest any) error {
	query := q.Build()

	// Rebind placeholders for the specific driver
	query = utils.RebindPlaceholder(query, q.executor.driver)

	fmt.Println("Executing:", query, "Args:", q.args)
	rows, err := q.executor.exec.Query(query, q.args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return utils.ScanRows(rows, dest)
}
