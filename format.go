package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func FormatAmount(a Amount, width int) string {
	redStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Width(width).
		Align(lipgloss.Right)
	greenStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Width(width).
		Align(lipgloss.Right)
	zeroStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Right)

	s := a.String()

	if a == 0 {
		return zeroStyle.Render(s)
	} else if a < 0 {
		return redStyle.Render(s)
	}
	return greenStyle.Render(s)
}

func Bold(s string) string {
	boldStyle := lipgloss.NewStyle().Bold(true)
	return boldStyle.Render(s)
}

func Bar(minus int, plus int, max int, width int) string {
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	minusWidth := width * minus / max
	plusWidth := width * plus / max
	return fmt.Sprintf(
		"%s%s : %s%s",
		strings.Repeat(" ", width-minusWidth),
		redStyle.Render(strings.Repeat("-", minusWidth)),
		greenStyle.Render(strings.Repeat("+", plusWidth)),
		strings.Repeat(" ", width-plusWidth),
	)
}
