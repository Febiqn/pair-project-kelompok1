package main

import "pair-project-kelompok1/handler"

func main() {
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
