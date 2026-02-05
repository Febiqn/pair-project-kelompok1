package config

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func RunMigration(db *sql.DB) {
	content, err := os.ReadFile("sql/ddl.sql")
	if err != nil {
		panic("Failed to read sql/ddl.sql")
	}

	queries := strings.Split(string(content), ";")

	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}

		_, err := db.Exec(q)
		if err != nil {
			panic("Failed to execute query: " + err.Error())
		}
	}

	fmt.Println("✔ Database schema migrated")
}
