package main

import (
	"cmp"
	"slices"
	"time"
)

type MonthlyExpenses struct {
	month       time.Month
	gains       Amount
	losses      Amount
	biggestGain *Expense
	biggestLoss *Expense
}

type YearlyExpenses struct {
	maxAmount Amount
	months    []MonthlyExpenses
}

type CategoryTotal struct {
	category string
	amount   Amount
}

func FilterByMonth(expenses []Expense, year int, month time.Month) []Expense {
	filtered := []Expense{}
	for _, expense := range expenses {
		if expense.Time.Year() == year && expense.Time.Month() == month {
			filtered = append(filtered, expense)
		}
	}

	return filtered
}

func GetCategoryTotals(expenses []Expense) []CategoryTotal {
	totalsMap := make(map[string]Amount)
	totals := []CategoryTotal{}

	for _, expense := range expenses {
		totalsMap[expense.Category()] += expense.Cost
	}

	for key, value := range totalsMap {
		totals = append(totals, CategoryTotal{
			category: key,
			amount:   value,
		})
	}

	slices.SortFunc(totals, func(a, b CategoryTotal) int {
		return cmp.Compare(a.amount, b.amount)
	})

	return totals
}

func GetMonthlyExpenses(expenses []Expense, year int, month time.Month) MonthlyExpenses {
	gains := 0
	losses := 0
	var biggestGain, biggestLoss *Expense

	filtered := FilterByMonth(expenses, year, month)

	for i, expense := range filtered {
		if expense.Cost < 0 {
			losses += -int(expense.Cost)
			if biggestLoss == nil || expense.Cost < biggestLoss.Cost {
				biggestLoss = &filtered[i]
			}
		} else {
			gains += int(expense.Cost)
			if biggestGain == nil || expense.Cost > biggestGain.Cost {
				biggestGain = &filtered[i]
			}
		}
	}

	return MonthlyExpenses{
		month:       month,
		gains:       Amount(gains),
		losses:      Amount(losses),
		biggestGain: biggestGain,
		biggestLoss: biggestLoss,
	}
}

func GetYearlyExpenses(expenses []Expense, year int) YearlyExpenses {
	months := []MonthlyExpenses{}
	maxAmount := Amount(1)

	for i := 1; i <= 12; i += 1 {
		monthly := GetMonthlyExpenses(expenses, year, time.Month(i))
		if monthly.gains > maxAmount {
			maxAmount = monthly.gains
		}
		if monthly.losses > maxAmount {
			maxAmount = monthly.losses
		}
		months = append(months, monthly)
	}

	return YearlyExpenses{
		months:    months,
		maxAmount: maxAmount,
	}
}
