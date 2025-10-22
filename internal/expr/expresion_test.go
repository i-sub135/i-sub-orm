package expr_test

import (
	"fmt"
	"testing"

	"github.com/i-sub135/i-sub-orm/internal/expr"
)

func TestCompile_Eq(t *testing.T) {
	tests := []struct {
		name     string
		input    expr.Eq
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "single field equality",
			input:    expr.Eq{"name": "John"},
			wantSQL:  "name = ?",
			wantArgs: []any{"John"},
		},
		{
			name:     "multiple fields equality",
			input:    expr.Eq{"name": "John", "age": 30},
			wantSQL:  "", // Will check both fields are present
			wantArgs: []any{"John", 30},
		},
		{
			name:     "integer value",
			input:    expr.Eq{"id": 1},
			wantSQL:  "id = ?",
			wantArgs: []any{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := expr.Compile(tt.input)

			fmt.Println("Generated SQL:", sql)
			fmt.Println("Arguments:", args)
			// For multiple fields, just check all args are present
			if len(tt.input) > 1 {
				if len(args) != len(tt.wantArgs) {
					t.Errorf("Compile() args count = %v, want %v", len(args), len(tt.wantArgs))
				}
				// Check SQL contains all field names and operators
				for k := range tt.input {
					if !contains(sql, k+" = ?") {
						t.Errorf("Compile() sql = %v, should contain '%s = ?'", sql, k)
					}
				}
			} else {
				if sql != tt.wantSQL {
					t.Errorf("Compile() sql = %v, want %v", sql, tt.wantSQL)
				}
				if !equalArgs(args, tt.wantArgs) {
					t.Errorf("Compile() args = %v, want %v", args, tt.wantArgs)
				}
			}
		})
	}
}

func TestCompile_Neq(t *testing.T) {
	tests := []struct {
		name     string
		input    expr.Neq
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "single field not equal",
			input:    expr.Neq{"status": "inactive"},
			wantSQL:  "status != ?",
			wantArgs: []any{"inactive"},
		},
		{
			name:     "multiple fields not equal",
			input:    expr.Neq{"status": "inactive", "deleted": true},
			wantSQL:  "",
			wantArgs: []any{"inactive", true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := expr.Compile(tt.input)

			if len(tt.input) > 1 {
				if len(args) != len(tt.wantArgs) {
					t.Errorf("Compile() args count = %v, want %v", len(args), len(tt.wantArgs))
				}
				for k := range tt.input {
					if !contains(sql, k+" != ?") {
						t.Errorf("Compile() sql = %v, should contain '%s != ?'", sql, k)
					}
				}
			} else {
				if sql != tt.wantSQL {
					t.Errorf("Compile() sql = %v, want %v", sql, tt.wantSQL)
				}
				if !equalArgs(args, tt.wantArgs) {
					t.Errorf("Compile() args = %v, want %v", args, tt.wantArgs)
				}
			}
		})
	}
}

func TestCompile_Gt(t *testing.T) {
	tests := []struct {
		name     string
		input    expr.Gt
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "greater than integer",
			input:    expr.Gt{"age": 18},
			wantSQL:  "age > ?",
			wantArgs: []any{18},
		},
		{
			name:     "greater than float",
			input:    expr.Gt{"price": 99.99},
			wantSQL:  "price > ?",
			wantArgs: []any{99.99},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := expr.Compile(tt.input)

			if sql != tt.wantSQL {
				t.Errorf("Compile() sql = %v, want %v", sql, tt.wantSQL)
			}
			if !equalArgs(args, tt.wantArgs) {
				t.Errorf("Compile() args = %v, want %v", args, tt.wantArgs)
			}
		})
	}
}

func TestCompile_Lt(t *testing.T) {
	tests := []struct {
		name     string
		input    expr.Lt
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "less than integer",
			input:    expr.Lt{"age": 65},
			wantSQL:  "age < ?",
			wantArgs: []any{65},
		},
		{
			name:     "less than float",
			input:    expr.Lt{"score": 50.5},
			wantSQL:  "score < ?",
			wantArgs: []any{50.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := expr.Compile(tt.input)

			if sql != tt.wantSQL {
				t.Errorf("Compile() sql = %v, want %v", sql, tt.wantSQL)
			}
			if !equalArgs(args, tt.wantArgs) {
				t.Errorf("Compile() args = %v, want %v", args, tt.wantArgs)
			}
		})
	}
}

func TestCompile_In(t *testing.T) {
	tests := []struct {
		name     string
		input    expr.In
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "IN with multiple values",
			input:    expr.In{"status": []any{"active", "pending", "review"}},
			wantSQL:  "status IN (?,?,?)",
			wantArgs: []any{"active", "pending", "review"},
		},
		{
			name:     "IN with single value",
			input:    expr.In{"id": []any{1}},
			wantSQL:  "id IN (?)",
			wantArgs: []any{1},
		},
		{
			name:     "IN with integer values",
			input:    expr.In{"id": []any{1, 2, 3, 4, 5}},
			wantSQL:  "id IN (?,?,?,?,?)",
			wantArgs: []any{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := expr.Compile(tt.input)

			if sql != tt.wantSQL {
				t.Errorf("Compile() sql = %v, want %v", sql, tt.wantSQL)
			}
			if !equalArgs(args, tt.wantArgs) {
				t.Errorf("Compile() args = %v, want %v", args, tt.wantArgs)
			}
		})
	}
}

func TestCompile_UnknownType(t *testing.T) {
	t.Run("unknown expression type returns empty", func(t *testing.T) {
		sql, args := expr.Compile("invalid")

		if sql != "" {
			t.Errorf("Compile() sql = %v, want empty string", sql)
		}
		if args != nil {
			t.Errorf("Compile() args = %v, want nil", args)
		}
	})
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func equalArgs(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
