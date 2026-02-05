package entity

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type ViewRevenue struct {
	PlaystationName string
	TotalBooking    int
	TotalRevenue    float64
}

func PrintRevenue(ViewRevenue []ViewRevenue) {
	data := [][]string{
		{"Name", "Total Booking", "Total Revenue"},
	}

	for _, v := range ViewRevenue {
		addData := []string{
			v.PlaystationName,
			fmt.Sprintf("%d", v.TotalBooking),
			fmt.Sprintf("%.2f", v.TotalRevenue),
		}
		data = append(data, addData)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	_ = table.Bulk(data[1:])
	_ = table.Render()
}

type ReportPS struct {
	PlaystationName string
	Condition       string
}

func PrintReportPS(ReportPS []ReportPS) {
	data := [][]string{
		{"Name", "Condition"},
	}

	for _, v := range ReportPS {
		addData := []string{
			v.PlaystationName, v.Condition}
		data = append(data, addData)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	_ = table.Bulk(data[1:])
	_ = table.Render()
}
