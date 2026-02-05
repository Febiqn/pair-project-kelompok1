package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"pair-project-kelompok1/entity"
	"strings"
)

func AdminFlow() {
	for {
		choice := ShowAdminMenu()

		switch choice {
		case "Update User Membership":
			updateMembership(db)
		case "View Revenue":
			showRevenue()
		case "Report Broken PS":
			ProcessReportAndFix(db)
		case "View PS Condition":
			ShowPSCondition()
		case "Exit":
			return
		}
	}
}

// Update membership
func UpdateMembershipQuery(db *sql.DB, memberStatus string, targetMember string) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	// Menggunakan LOWER untuk pencarian case-insensitive agar lebih user-friendly
	query := `UPDATE users SET membership_status = ? WHERE LOWER(TRIM(name)) = LOWER(?)`
	result, err := db.Exec(query, memberStatus, targetMember)
	if err != nil {
		return 0, fmt.Errorf("query execution failed: %w", err)
	}

	return result.RowsAffected()
}

func updateMembership(db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("\nEnter the name of Membership you want to update: ")
	if !scanner.Scan() {
		return
	}
	targetMember := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter status of membership (ACTIVE/INACTIVE): ")
	if !scanner.Scan() {
		return
	}
	memberStatus := strings.ToUpper(strings.TrimSpace(scanner.Text()))

	if targetMember == "" || memberStatus == "" {
		fmt.Println("\n[!] Error: Name and status cannot be empty.")
		return
	}

	// Menangkap rowsAffected untuk memberi feedback ke user
	rows, err := UpdateMembershipQuery(db, memberStatus, targetMember)
	if err != nil {
		fmt.Printf("[!] Update failed: %v\n", err)
		return
	}

	if rows == 0 {
		fmt.Printf("[?] No member found with name: %s\n", targetMember)
	} else {
		fmt.Printf("[+] Success! %d record(s) updated.\n", rows)
	}
}

// showing revenue
func showRevenue() {
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	query := `
		SELECT 
    COALESCE(p.ps_name, 'TOTAL >') AS ps_name, 
    COUNT(b.bill_id) AS total_trx, 
    COALESCE(SUM(b.total_amount), 0) AS total_revenue
	FROM playstations p
	LEFT JOIN rentals r ON p.ps_id = r.ps_id
	LEFT JOIN billing b ON r.rental_id = b.rental_id
	GROUP BY p.ps_name WITH ROLLUP;
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	var reports []entity.ViewRevenue
	for rows.Next() {
		var r entity.ViewRevenue
		err := rows.Scan(&r.PlaystationName, &r.TotalBooking, &r.TotalRevenue)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		reports = append(reports, r)
	}

	fmt.Println("\nREVENUE REPORT")
	entity.PrintRevenue(reports)
}

// report condition PS
// retrive data ftom DB and
func FetchAllPlaystations(db *sql.DB) ([]entity.ReportPS, error) {
	query := `SELECT ps_name, condition_status FROM playstations`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []entity.ReportPS
	for rows.Next() {
		var r entity.ReportPS
		if err := rows.Scan(&r.PlaystationName, &r.Condition); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	fmt.Println("\n    Latest Data Playstation    ")
	if len(reports) == 0 {
		fmt.Println("Database is empty.")
	} else {
		// Menggunakan helper function yang kamu buat
		entity.PrintReportPS(reports)
	}

	return reports, nil
}

// change status based on name
func UpdateCondition(db *sql.DB, psName string, newCondition string) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("Database not initialized")
	}

	query := `UPDATE playstations SET condition_status = ? WHERE TRIM(ps_name) = ?`
	result, err := db.Exec(query, newCondition, psName)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Process Report
func ProcessReportAndFix(db *sql.DB) {
	_, err := FetchAllPlaystations(db)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}

	// Gunakan bufio scanner agar bisa membaca spasi dan membersihkan newline (\r\n)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the name of the PS you want to update.: ")
	scanner.Scan()
	targetPS := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter New Condition (eg: AVAILABLE/BROKEN)): ")
	scanner.Scan()
	newCondition := strings.TrimSpace(scanner.Text())

	// Validasi input kosong
	if targetPS == "" || newCondition == "" {
		fmt.Println("\n[!] Error: PS Name and Condition cannot be empty.")
		return
	}

	// 2. Eksekusi Update
	_, err = UpdateCondition(db, targetPS, newCondition)
	if err != nil {
		fmt.Println("Update failed:", err)
		return
	}
}

// show PS condition
func ShowPSCondition() {
	if db == nil {
		fmt.Println("Database not initialized")
		return
	}

	query := `
	SELECT ps_name, condition_status FROM playstations;
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	var reports []entity.ConditionPS
	for rows.Next() {
		var r entity.ConditionPS
		err := rows.Scan(&r.PlaystationName, &r.Condition)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		reports = append(reports, r)
	}

	entity.PrintViewPSCondition(reports)
}
