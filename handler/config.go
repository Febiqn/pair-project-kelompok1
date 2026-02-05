package handler

import "database/sql"

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}
