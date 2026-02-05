package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFlag := flag.Bool("csv", false, "Export expenses as CSV to stdout")
	flag.Parse()


	if flag.NArg() < 1 {
		fmt.Println("No filename provided.")
		os.Exit(1)
	}
	
	filename := flag.Arg(0)

	expenses, error := ReadExpenses(filename)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	if *csvFlag {
		for _, expense := range expenses {
			fmt.Println(expense.ToCSV())
		}
		return
	}

	year, month, _ := time.Now().Date()
	monthly := GetMonthlyExpenses(expenses, year, month)

	fmt.Println(Bold("This month's expenses:"))
	fmt.Printf("Gains:  %s\n", FormatAmount(monthly.gains, 10))
	fmt.Printf("Losses: %s\n", FormatAmount(-monthly.losses, 10))
	fmt.Println()
	fmt.Printf("Total:  %s\n", FormatAmount(monthly.gains-monthly.losses, 10))
	fmt.Println()
	if monthly.biggestGain != nil {
		fmt.Printf("Biggest gain: %s (%s)\n", FormatAmount(monthly.biggestGain.Cost, 0), monthly.biggestGain.Description)
	}
	if monthly.biggestLoss != nil {
		fmt.Printf("Biggest loss: %s (%s)\n", FormatAmount(monthly.biggestLoss.Cost, 0), monthly.biggestLoss.Description)
	}
	fmt.Println()

	fmt.Println(Bold("By Category:"))

	for _, total := range GetCategoryTotals(FilterByMonth(expenses, year, month)) {
		name := total.category
		if name == "" {
			name = "N/A"
		}
		fmt.Printf("%-10s %s\n", name, FormatAmount(total.amount, 10))
	}
	fmt.Println()

	yearlyExpenses := GetYearlyExpenses(expenses, year)

	var yearlyReport strings.Builder

	yearlyReport.WriteString(Bold(fmt.Sprintf("Yearly report (%d)", year)) + "\n")
	width := 10
	var total Amount
	for i, monthly := range yearlyExpenses.months {
		subTotal := monthly.gains - monthly.losses
		total += subTotal
		fmt.Fprintf(
			&yearlyReport,
			"%5s   %s %s\n",
			time.Month(i + 1).String()[:3],
			Bar(int(monthly.losses), int(monthly.gains), int(yearlyExpenses.maxAmount), width),
			FormatAmount(subTotal, 10),
		)
	}

	fmt.Println(yearlyReport.String())
	fmt.Printf("Total   %s %s\n", strings.Repeat(" ", 2*width+3), FormatAmount(total, 10))
}
