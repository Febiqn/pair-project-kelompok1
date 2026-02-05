package handler

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

	// ❌ If user does NOT want membership → DO NOTHING
	if choice == "No" {
		fmt.Println("Membership not activated. Registration cancelled.")
		return
	}

	// ✅ Only reach here if user chose YES
	status := "ACTIVE"
	memberNo := generateMembershipNumber()

	query := `
		INSERT INTO users (user_name, membership_status, membership_number)
		VALUES (?, ?, ?)
	`

	_, err = db.Exec(query, name, status, memberNo)
	if err != nil {
		fmt.Println("Failed to register user:", err)
		return
	}

	fmt.Println("✔ Membership registered successfully")
	fmt.Println("✔ Membership number:", memberNo)
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
		WHERE user_name = ? AND membership_number = ?
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
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	rolePrompt := promptui.Select{
		Label: "Are you a member?",
		Items: []string{"Yes", "No"},
	}

	_, role, err := rolePrompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	var (
		userID           int
		userName         string
		membershipStatus string
	)

	if role == "Yes" {
		memberPrompt := promptui.Prompt{
			Label: "Enter membership number (PS-??)",
		}

		memberNo, err := memberPrompt.Run()
		if err != nil {
			fmt.Println("Cancelled")
			return
		}

		err = db.QueryRow(`
			SELECT user_id, user_name, membership_status
			FROM users
			WHERE membership_number = ?
			  AND membership_status = 'ACTIVE'
		`, memberNo).Scan(&userID, &userName, &membershipStatus)

		if err == sql.ErrNoRows {
			fmt.Println("Active member not found")
			return
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	if role == "No" {
		name, err := promptName()
		if err != nil {
			fmt.Println("Cancelled")
			return
		}

		res, err := db.Exec(`
			INSERT INTO users (user_name, membership_status)
			VALUES (?, 'INACTIVE')
		`, name)

		if err != nil {
			fmt.Println("Failed to register non-member:", err)
			return
		}

		id, _ := res.LastInsertId()
		userID = int(id)
		userName = name
		membershipStatus = "INACTIVE"
	}

	rows, err := db.Query(`
		SELECT ps_id, ps_name
		FROM playstations
		WHERE condition_status = 'AVAILABLE'
	`)
	if err != nil {
		fmt.Println("Failed to fetch PlayStations:", err)
		return
	}
	defer rows.Close()

	type PS struct {
		ID   int
		Name string
	}

	var psList []PS
	var psNames []string

	for rows.Next() {
		var ps PS
		rows.Scan(&ps.ID, &ps.Name)
		psList = append(psList, ps)
		psNames = append(psNames, ps.Name)
	}

	if len(psList) == 0 {
		fmt.Println("No PlayStation available")
		return
	}

	psPrompt := promptui.Select{
		Label: "Select PlayStation",
		Items: psNames,
	}

	index, _, err := psPrompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	psID := psList[index].ID
	psName := psList[index].Name

	durationPrompt := promptui.Prompt{
		Label: "Enter rental duration (hours)",
		Validate: func(input string) error {
			val, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil || val <= 0 {
				return fmt.Errorf("Enter a valid number")
			}
			return nil
		},
	}

	durationStr, err := durationPrompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	duration, _ := strconv.Atoi(durationStr)

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Hour)
	baseAmount := float64(duration * 10000)

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Transaction failed:", err)
		return
	}

	// INSERT RENTAL
	res, err := tx.Exec(`
		INSERT INTO rentals (
			user_id,
			user_name,
			ps_id,
			ps_name,
			membership_status,
			start_time,
			duration_hours,
			end_time
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, userName, psID, psName, membershipStatus, startTime, duration, endTime)

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to create rental:", err)
		return
	}

	rentalID, _ := res.LastInsertId()

	_, err = tx.Exec(`
		UPDATE playstations
		SET condition_status = 'RENTED'
		WHERE ps_id = ?
	`, psID)

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to update PlayStation:", err)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO billing (rental_id, total_amount)
		VALUES (?, ?)
	`, rentalID, baseAmount)

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to create billing:", err)
		return
	}

	tx.Commit()

	fmt.Println("\n✔ Rental successful")
	fmt.Println("User       :", userName)
	fmt.Println("Membership :", membershipStatus)
	fmt.Println("PlayStation:", psName)
	fmt.Println("Start      :", startTime.Format("2006-01-02 15:04"))
	fmt.Println("End        :", endTime.Format("2006-01-02 15:04"))
	fmt.Println("Base Price : Rp", baseAmount)
}

func checkTime() {
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	rows, err := db.Query(`
		SELECT
			r.rental_id,
			r.ps_name,
			r.end_time
		FROM rentals r
		WHERE r.status = 'ONGOING'
	`)
	if err != nil {
		fmt.Println("Failed to fetch ongoing rentals:", err)
		return
	}
	defer rows.Close()

	type Rental struct {
		ID      int
		PsName  string
		EndTime time.Time
	}

	var rentals []Rental
	var psNames []string

	for rows.Next() {
		var r Rental
		rows.Scan(&r.ID, &r.PsName, &r.EndTime)
		rentals = append(rentals, r)
		psNames = append(psNames, r.PsName)
	}

	if len(rentals) == 0 {
		fmt.Println("No ongoing rentals")
		return
	}

	psPrompt := promptui.Select{
		Label: "Select PlayStation",
		Items: psNames,
	}

	index, _, err := psPrompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	selected := rentals[index]

	// ========================
	// CALCULATE TIME LEFT
	// ========================
	remaining := time.Until(selected.EndTime)

	if remaining <= 0 {
		fmt.Println("⏰ Rental time has ended")
		return
	}

	minutesLeft := int(remaining.Minutes())

	// ========================
	// DISPLAY
	// ========================
	fmt.Println("\n⏳ Time Remaining")
	fmt.Printf("%d minutes\n", minutesLeft)
}

func payBill() {
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	// ========================
	// FETCH UNPAID RENTALS
	// ========================
	rows, err := db.Query(`
		SELECT
			r.rental_id,
			r.user_name,
			r.ps_id,
			r.ps_name,
			r.duration_hours,
			r.membership_status
		FROM rentals r
		JOIN billing b ON r.rental_id = b.rental_id
		WHERE r.status = 'ONGOING'
		  AND b.bill_status = 'UNPAID'
	`)
	if err != nil {
		fmt.Println("Failed to fetch unpaid rentals:", err)
		return
	}
	defer rows.Close()

	type Rental struct {
		ID       int
		UserName string
		PsID     int
		PsName   string
		Duration int
		Member   string
	}

	var rentals []Rental
	var items []string

	for rows.Next() {
		var r Rental
		rows.Scan(
			&r.ID,
			&r.UserName,
			&r.PsID,
			&r.PsName,
			&r.Duration,
			&r.Member,
		)

		rentals = append(rentals, r)
		items = append(items, fmt.Sprintf(
			"%s - %s (%d hrs)",
			r.UserName,
			r.PsName,
			r.Duration,
		))
	}

	if len(rentals) == 0 {
		fmt.Println("❌ No unpaid rentals")
		return
	}

	// ========================
	// SELECT RENTAL
	// ========================
	selectPrompt := promptui.Select{
		Label: "Select rental to pay",
		Items: items,
	}

	index, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	selected := rentals[index]

	// ========================
	// CALCULATE BILL
	// ========================
	baseAmount := float64(selected.Duration * 10000)
	discount := 0.0

	if selected.Member == "ACTIVE" {
		discount = baseAmount * 0.10
	}

	finalAmount := baseAmount - discount

	// ========================
	// SHOW BILL SUMMARY
	// ========================
	fmt.Println("\nTOTAL BILL")
	fmt.Println("User        :", selected.UserName)
	fmt.Println("PlayStation :", selected.PsName)
	fmt.Println("Duration    :", selected.Duration, "hours")
	fmt.Println("Base Price  : Rp", baseAmount)
	fmt.Println("Discount    : Rp", discount)
	fmt.Println("Total Pay   : Rp", finalAmount)

	confirmPrompt := promptui.Select{
		Label: "Confirm payment?",
		Items: []string{"Yes", "Cancel"},
	}

	_, confirm, _ := confirmPrompt.Run()
	if confirm == "Cancel" {
		fmt.Println("Payment cancelled")
		return
	}

	// ========================
	// TRANSACTION
	// ========================
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Transaction failed:", err)
		return
	}

	// UPDATE BILLING
	_, err = tx.Exec(`
		UPDATE billing
		SET total_amount = ?,
		    bill_status = 'PAID',
		    paid_at = NOW()
		WHERE rental_id = ?
	`, finalAmount, selected.ID)

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to update billing:", err)
		return
	}

	// COMPLETE RENTAL
	_, err = tx.Exec(`
		UPDATE rentals
		SET status = 'COMPLETED'
		WHERE rental_id = ?
	`, selected.ID)

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to complete rental:", err)
		return
	}

	// RELEASE PLAYSTATION
	_, err = tx.Exec(`
		UPDATE playstations
		SET condition_status = 'AVAILABLE'
		WHERE ps_id = ?
	`, selected.PsID)

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to release PlayStation:", err)
		return
	}

	tx.Commit()

	fmt.Println("\n✔ Payment successful")
	fmt.Println("✔ Rental completed")
	fmt.Println("✔ PlayStation is now AVAILABLE")
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
