package main

import (
	"log"
	"time"

	"github.com/i-sub135/i-sub-orm/internal/driver"
	// "github.com/i-sub135/i-sub-orm/internal/expr"
	"github.com/i-sub135/i-sub-orm/pkg/orm"
	_ "github.com/lib/pq" // Import postgres driver
	// _ "github.com/mattn/go-sqlite3" // Import sqlite3 driver
)

// User represents the users table in the database
type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// DATABASE_URL="postgresql://tracking_user:tracking_pass@localhost:5432/tracking_db"
func main() {
	db, err := orm.Open(driver.Postgres.String(), "host=localhost port=5432 user=tracking_user password=tracking_pass dbname=tracking_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Example: Get all users
	var users []User
	err = db.Table("users").
		Select("*").
		Get(&users)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d users", len(users))
	log.Printf("%+v", users)
}
