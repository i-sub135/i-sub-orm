package utils_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/i-sub135/i-sub-orm/internal/constant"
	"github.com/i-sub135/i-sub-orm/internal/utils"
)

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

func TestScanRows_IntoSlice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "age"}).
		AddRow(1, "John Doe", "john@example.com", 25).
		AddRow(2, "Jane Smith", "jane@example.com", 30).
		AddRow(3, "Bob Wilson", "bob@example.com", 35)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var users []User
	err = utils.ScanRows(queryRows, &users)
	if err != nil {
		t.Fatalf("ScanRows failed: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}

	// Verify first user
	if users[0].ID != 1 || users[0].Name != "John Doe" || users[0].Email != "john@example.com" || users[0].Age != 25 {
		t.Errorf("first user data mismatch: %+v", users[0])
	}

	// Verify second user
	if users[1].ID != 2 || users[1].Name != "Jane Smith" || users[1].Email != "jane@example.com" || users[1].Age != 30 {
		t.Errorf("second user data mismatch: %+v", users[1])
	}

	// Verify third user
	if users[2].ID != 3 || users[2].Name != "Bob Wilson" || users[2].Email != "bob@example.com" || users[2].Age != 35 {
		t.Errorf("third user data mismatch: %+v", users[2])
	}
}

func TestScanRows_IntoStruct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "age"}).
		AddRow(1, "John Doe", "john@example.com", 25)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email, age FROM users WHERE id = 1")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var user User
	err = utils.ScanRows(queryRows, &user)
	if err != nil {
		t.Fatalf("ScanRows failed: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got '%s'", user.Email)
	}
	if user.Age != 25 {
		t.Errorf("expected age 25, got %d", user.Age)
	}
}

func TestScanRows_StructNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "age"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email, age FROM users WHERE id = 999")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var user User
	err = utils.ScanRows(queryRows, &user)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestScanRows_EmptySlice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "age"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email, age FROM users WHERE 1=0")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var users []User
	err = utils.ScanRows(queryRows, &users)
	if err != nil {
		t.Fatalf("ScanRows failed: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("expected empty slice, got %d items", len(users))
	}
}

func TestScanRows_WithoutDBTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Product A", 99.99)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var product Product
	err = utils.ScanRows(queryRows, &product)
	if err != nil {
		t.Fatalf("ScanRows failed: %v", err)
	}

	if product.ID != 1 || product.Name != "Product A" || product.Price != 99.99 {
		t.Errorf("product data mismatch: %+v", product)
	}
}

func TestScanRows_ErrorNotPointer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var user User
	err = utils.ScanRows(queryRows, user) // Not a pointer
	if err != constant.ErrDestination {
		t.Errorf("expected ErrDestination, got %v", err)
	}
}

func TestScanRows_ErrorInvalidType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var invalidDest string
	err = utils.ScanRows(queryRows, &invalidDest) // Invalid type (not slice or struct)
	if err != constant.ErrDestinationType {
		t.Errorf("expected ErrDestinationType, got %v", err)
	}
}

func TestScanRows_PartialColumns(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	// Query returns only 2 columns, but struct has 4 fields
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "John Doe")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var user User
	err = utils.ScanRows(queryRows, &user)
	if err != nil {
		t.Fatalf("ScanRows failed: %v", err)
	}

	// Only ID and Name should be set
	if user.ID != 1 || user.Name != "John Doe" {
		t.Errorf("expected ID=1 and Name='John Doe', got %+v", user)
	}
	// Email and Age should be zero values
	if user.Email != "" || user.Age != 0 {
		t.Errorf("expected zero values for Email and Age, got %+v", user)
	}
}

func TestScanRows_ExtraColumns(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	// Query returns extra columns that don't exist in struct
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "extra_col"}).
		AddRow(1, "John Doe", "john@example.com", 25, "ignored")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email, age, extra_col FROM users")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer queryRows.Close()

	var user User
	err = utils.ScanRows(queryRows, &user)
	if err != nil {
		t.Fatalf("ScanRows failed: %v", err)
	}

	// All struct fields should be populated, extra column ignored
	if user.ID != 1 || user.Name != "John Doe" || user.Email != "john@example.com" || user.Age != 25 {
		t.Errorf("user data mismatch: %+v", user)
	}
}
