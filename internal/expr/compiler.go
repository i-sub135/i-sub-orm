package expr

import (
	"fmt"
	"strings"
)

// compile will compile the given condition into a SQL string and its arguments.
func Compile(condition any) (string, []any) {
	switch cond := condition.(type) {
	case Eq:
		return builCompair(cond, "=")
	case Neq:
		return builCompair(cond, "!=")
	case Gt:
		return builCompair(cond, ">")
	case Lt:
		return builCompair(cond, "<")
	case In:
		return buildIN(cond)
	default:
		return "", nil
	}

}

// builCompair builds comparison expressions like "field = ?" and returns the SQL string and arguments.
func builCompair(data map[string]any, operator string) (string, []any) {

	parts := make([]string, 0, len(data))
	args := make([]any, 0, len(data))
	for k, v := range data {
		parts = append(parts, fmt.Sprintf("%s %s ?", k, operator))
		args = append(args, v)
	}

	return strings.Join(parts, " AND "), args

}

// buildIN builds IN expressions like "field IN (?, ?, ?)" and returns the SQL string and arguments.
func buildIN(data map[string][]any) (string, []any) {
	parts := make([]string, 0, len(data))
	args := make([]any, 0)

	for k, v := range data {
		placeholders := strings.Repeat("?,", len(v))
		placeholders = strings.TrimRight(placeholders, ",")
		parts = append(parts, fmt.Sprintf("%s IN (%s)", k, placeholders))
		args = append(args, v...)
	}
	return strings.Join(parts, " AND "), args
}
