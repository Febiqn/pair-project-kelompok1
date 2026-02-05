package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() {

	db, err := sql.Open("mysql", "appuser:app123@tcp(localhost:3306)/rental?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Success connect to DB")
}
