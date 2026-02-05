package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Format
// 2026-01-10 -$453.89 #gift Gift card
// 2026-01-11 +$2000   #job  Salary

type Amount int

func (a Amount) String() string {
	value := int(a)
	if value < 0 {
		value *= -1
	}

	cents := value % 100
	units := value / 100
	if int(a) < 0 {
		return fmt.Sprintf("-$%v.%02v", units, cents)
	}
	return fmt.Sprintf("$%v.%02v", units, cents)
}

type Expense struct {
	Time        time.Time
	Cost        Amount
	Description string
	Tags        []string
}

var (
	ErrInvalidFormat = errors.New("Invalid expense format")
	ErrInvalidDate   = errors.New("Invalid date")
)

var (
	expenseRegexp  = regexp.MustCompile(`^.*(\d{4}-\d\d-\d\d)\s+([+-]?)\$(\d+(\.\d\d|))`)
	tagsRegexp     = regexp.MustCompile(`\s*#([^\s#]+)`)
)

func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

func (e *Expense) Category() string {
	if len(e.Tags) == 0 {
		return ""
	}
	return e.Tags[0]
}

func (e Expense) ToCSV() string {
	return fmt.Sprintf("%s,%s,%s,%s", e.Time.Format("2006-01-02"), e.Cost.String(), escapeCSV(strings.Join(e.Tags, " ")), escapeCSV(e.Description))
}

func ParseExpense(text string) (Expense, error) {
	matches := expenseRegexp.FindStringSubmatch(text)

	if matches == nil {
		return Expense{}, ErrInvalidFormat
	}

	dateString := matches[1]
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return Expense{}, fmt.Errorf("%w: %s", ErrInvalidDate, dateString)
	}

	cost, _ := strconv.Atoi(strings.ReplaceAll(matches[3], ".", ""))
	if matches[4] == "" {
		cost *= 100
	}
	if matches[2] == "-" {
		cost *= -1
	}

	text = text[len(matches[0]):]

	tagsMap := make(map[string]bool)
	tags := []string{}
	for _, match := range tagsRegexp.FindAllStringSubmatch(text, -1) {
		if !tagsMap[match[1]] {
			tags = append(tags, match[1])
		}
		tagsMap[match[1]] = true
		text = text[len(match[0]):]
	}

	expense := Expense{
		Time:        date,
		Cost:        Amount(cost),
		Description: strings.TrimSpace(text),
		Tags:        tags,
	}

	return expense, nil
}
