package handler

import (
	"fmt"
	"pair-project-kelompok1/entity"
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

func reportBroken() {
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

	var reports []entity.ReportPS
	for rows.Next() {
		var r entity.ReportPS
		err := rows.Scan(&r.PlaystationName, &r.Condition)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		reports = append(reports, r)
	}

	entity.PrintReportPS(reports)
}
