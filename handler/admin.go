package handler

import (
	"fmt"
)

func AdminFlow() {
	for {
		choice := ShowAdminMenu()

		switch choice {
		case "Update User Membership":
			updateMembership()
		case "View Revenue":
			showRevenue()
		case "Report Broken PS":
			reportBroken()
		case "Exit":
			return
		}
	}
}

func updateMembership() {
	fmt.Println("UPDATE MEMBER")
}

func showRevenue() {
	fmt.Println("Show")
}

func reportBroken() {
	fmt.Println("REPROT")
}
