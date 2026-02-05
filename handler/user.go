package handler

import "fmt"

func UserFlow() {
	for {
		choice := ShowUserMenu()

		switch choice {
		case "Register Membership":
			registerUser()
		case "Rent PlayStation":
			rentPS()
		case "Check Time Left":
			checkTime()
		case "Pay Bill":
			payBill()
		case "Exit":
			return
		}
	}
}

func registerUser() {
	fmt.Println("REG")
}

func rentPS() {
	fmt.Println("RENTE")
}

func checkTime() {
	fmt.Println("CHECK")
}

func payBill() {
	fmt.Println("PAY")
}
