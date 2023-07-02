package examples

import (
	"database/sql"
)

func AnaylzeTypesExample() {
	var db *sql.DB
	db.Exec("SELECT * FROM users")
}