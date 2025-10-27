package utils

import (
	"fmt"
	"strings"
)

// RebindPlaceholder converts ? placeholders to the format used by the specified driver
// For postgres: $1, $2, $3
// For mysql/sqlite: ? (no change)
func RebindPlaceholder(query string, driver string) string {
	switch driver {
	case "postgres", "postgresql":
		return rebindPostgres(query)
	default:
		return query
	}
}

// rebindPostgres converts ? to $1, $2, $3, etc.
func rebindPostgres(query string) string {
	count := 1
	var result strings.Builder

	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			result.WriteString(fmt.Sprintf("$%d", count))
			count++
		} else {
			result.WriteByte(query[i])
		}
	}

	return result.String()
}
