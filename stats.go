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
	topGains    []Expense
	topLosses   []Expense
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

func GetMonthlyExpenses(expenses []Expense, year int, month time.Month, topN int) MonthlyExpenses {
	gains := 0
	losses := 0
	var lossExpenses []Expense
	var gainExpenses []Expense

	filtered := FilterByMonth(expenses, year, month)

	for _, expense := range filtered {
		if expense.Cost < 0 {
			losses += -int(expense.Cost)
			lossExpenses = append(lossExpenses, expense)
		} else {
			gains += int(expense.Cost)
			gainExpenses = append(gainExpenses, expense)
		}
	}

	slices.SortFunc(lossExpenses, func(a, b Expense) int {
		return cmp.Compare(a.Cost, b.Cost)
	})
	slices.SortFunc(gainExpenses, func(a, b Expense) int {
		return cmp.Compare(a.Cost, b.Cost)
	})
	if topN > 0 && len(lossExpenses) > topN {
		lossExpenses = lossExpenses[:topN]
	}
	if topN > 0 && len(gainExpenses) > topN {
		gainExpenses = gainExpenses[:topN]
	}

	return MonthlyExpenses{
		month:       month,
		gains:       Amount(gains),
		losses:      Amount(losses),
		topGains:    gainExpenses,
		topLosses:   lossExpenses,
	}
}

func GetYearlyExpenses(expenses []Expense, year int) YearlyExpenses {
	months := []MonthlyExpenses{}
	maxAmount := Amount(1)

	for i := 1; i <= 12; i += 1 {
		monthly := GetMonthlyExpenses(expenses, year, time.Month(i), 0)
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
