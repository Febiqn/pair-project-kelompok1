package handler

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/manifoldco/promptui"
)

func UserFlow() {
	for {
		choice := ShowUserMenu()

		switch choice {
		case "Register Membership":
			registerUser()
		case "Check Membership":
			checkMember()
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
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	name, err := promptName()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	activatePrompt := promptui.Select{
		Label: "Activate membership now?",
		Items: []string{"Yes", "No"},
	}

	_, choice, err := activatePrompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	status := "INACTIVE"
	var memberNo *string

	if choice == "Yes" {
		status = "ACTIVE"
		no := generateMembershipNumber()
		memberNo = &no
	}

	query := `
		INSERT INTO users (name, membership_status, membership_number)
		VALUES (?, ?, ?)
	`

	_, err = db.Exec(query, name, status, memberNo)
	if err != nil {
		fmt.Println("Failed to register user:", err)
		return
	}

	fmt.Println("✔ User registered successfully")
	fmt.Println("✔ Membership status:", status)
	if memberNo != nil {
		fmt.Println("✔ Membership number:", *memberNo)
	}
}

func checkMember() {
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	name, err := promptName()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	memberPrompt := promptui.Prompt{
		Label: "Enter membership number (PS-??)",
	}

	memberNo, err := memberPrompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	var userID int
	var status string

	query := `
		SELECT user_id, membership_status
		FROM users
		WHERE name = ? AND membership_number = ?
	`

	err = db.QueryRow(query, name, memberNo).Scan(&userID, &status)
	if err == sql.ErrNoRows {
		fmt.Println("Member not found")
		return
	} else if err != nil {
		fmt.Println("Error checking member:", err)
		return
	}

	fmt.Println("\nMember found")
	fmt.Println("Name:", name)
	fmt.Println("Membership Number:", memberNo)
	fmt.Println("Status:", status)

	actionPrompt := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{"Delete Membership", "Back to Menu"},
	}

	_, action, _ := actionPrompt.Run()

	if action == "Delete Membership" {
		_, err := db.Exec(`DELETE FROM users WHERE user_id = ?`, userID)
		if err != nil {
			fmt.Println("Failed to delete membership:", err)
			return
		}
		fmt.Println("Membership deleted successfully")
	}
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

func generateMembershipNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("PS-%d", rand.Intn(99)+1)
}

func promptName() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter your name",
		Validate: func(input string) error {
			if len(input) < 5 {
				return fmt.Errorf("Enter 5 characters name")
			}
			return nil
		},
	}

	return prompt.Run()
}
