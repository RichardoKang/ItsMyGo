package main

import (
	"database/sql"
)

func main() {
	db, _ := sql.Open("sqlite3", "gee.db")

	defer func() {
		_ = db.Close()
	}()

}
