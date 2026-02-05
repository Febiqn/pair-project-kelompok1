package main

import (
	"pair-project-kelompok1/config"
	"pair-project-kelompok1/handler"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	config.RunMigration(db)

	handler.InitDB(db)

	for {
		role := handler.RoleMenu()

		switch role {
		case "User":
			handler.UserFlow()
		case "Admin":
			handler.AdminFlow()
		case "Exit":
			return
		}
	}
}
