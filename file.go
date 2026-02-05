package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

func ReadExpenses(filename string) ([]Expense, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	expenses := []Expense{}

	for i, line := range lines {
		expense, err := ParseExpense(line)
		if errors.Is(err, ErrInvalidDate) {
			return expenses, fmt.Errorf("Line %d: %w", i+1, err)
		} else if err == nil {
			expenses = append(expenses, expense)
		}
	}

	slices.SortFunc(expenses, func(a, b Expense) int {
		return a.Time.Compare(b.Time)
	})

	return expenses, nil
}
