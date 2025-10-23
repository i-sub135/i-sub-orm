package main

import (
	"log"

	"github.com/i-sub135/i-sub-orm/internal/expr"
	"github.com/i-sub135/i-sub-orm/pkg/orm"
	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver
)

func main() {
	db, err := orm.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Example chaining with expression-style WHEREs
	err = db.Table("users").
		Select("*").
		Where(expr.Eq{"status": "active", "role": "admin"}).
		Where(expr.Gt{"age": 18}).
		Get(nil)

	if err != nil {
		log.Fatal(err)
	}
}
