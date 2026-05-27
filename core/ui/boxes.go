/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : boxes.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Helpers that render bordered, neon-colored boxes used
 *                across the CLI for titles, summaries, and notices.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// boxStyle returns a bordered lipgloss style with the given border color.
func boxStyle(borderColor lipgloss.Color, contentColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Foreground(contentColor).
		Padding(0, 2).
		MarginTop(1).
		MarginBottom(1)
}

// Box prints a neon-bordered information box.
func Box(title string, body string) {
	titleStyled := StyleTitle.Render(title)
	bodyStyled := lipgloss.NewStyle().Foreground(ColorWhite).Render(body)
	content := titleStyled + "\n" + strings.Repeat("─", lipgloss.Width(title)) + "\n" + bodyStyled
	fmt.Println(boxStyle(ColorBlue, ColorWhite).Render(content))
}

// SuccessBox highlights a positive milestone.
func SuccessBox(title string, body string) {
	titleStyled := StyleSuccess.Render("✓ " + title)
	bodyStyled := lipgloss.NewStyle().Foreground(ColorGreen).Render(body)
	content := titleStyled + "\n" + strings.Repeat("─", lipgloss.Width(title)+2) + "\n" + bodyStyled
	fmt.Println(boxStyle(ColorGreen, ColorGreen).Render(content))
}

// ErrorBox renders a fatal error in a red bordered box.
func ErrorBox(title string, body string) {
	titleStyled := StyleError.Render("✗ " + title)
	bodyStyled := lipgloss.NewStyle().Foreground(ColorRed).Render(body)
	content := titleStyled + "\n" + strings.Repeat("─", lipgloss.Width(title)+2) + "\n" + bodyStyled
	fmt.Println(boxStyle(ColorRed, ColorRed).Render(content))
}

// WarningBox renders a recoverable warning in a yellow bordered box.
func WarningBox(title string, body string) {
	titleStyled := StyleWarning.Render("! " + title)
	bodyStyled := lipgloss.NewStyle().Foreground(ColorYellow).Render(body)
	content := titleStyled + "\n" + strings.Repeat("─", lipgloss.Width(title)+2) + "\n" + bodyStyled
	fmt.Println(boxStyle(ColorYellow, ColorYellow).Render(content))
}

// KeyValueBox prints a labeled key/value summary inside a bordered box.
func KeyValueBox(title string, pairs [][2]string) {
	titleStyled := StyleTitle.Render(title)
	var lines []string
	lines = append(lines, titleStyled)
	lines = append(lines, strings.Repeat("─", lipgloss.Width(title)))
	for _, p := range pairs {
		lines = append(lines, StyleLabel.Render(p[0]+":")+" "+StyleValue.Render(p[1]))
	}
	fmt.Println(boxStyle(ColorPurple, ColorWhite).Render(strings.Join(lines, "\n")))
}

// SectionHeader renders a slim colored section header.
func SectionHeader(text string) {
	bar := strings.Repeat("═", lipgloss.Width(text)+4)
	fmt.Println(StyleTitle.Render("╔" + bar + "╗"))
	fmt.Println(StyleTitle.Render("║  " + text + "  ║"))
	fmt.Println(StyleTitle.Render("╚" + bar + "╝"))
}

// Info prints an informational line with a colored prefix.
func Info(msg string) {
	fmt.Println(StyleInfo.Render("[i] ") + msg)
}

// Success prints a success line with a colored prefix.
func Success(msg string) {
	fmt.Println(StyleSuccess.Render("[+] ") + msg)
}

// Warning prints a warning line with a colored prefix.
func Warning(msg string) {
	fmt.Println(StyleWarning.Render("[!] ") + msg)
}

// Error prints an error line with a colored prefix.
func Error(msg string) {
	fmt.Println(StyleError.Render("[x] ") + msg)
}

// Step prints a numbered step entry in the workflow.
func Step(idx int, total int, msg string) {
	prefix := fmt.Sprintf("[%d/%d]", idx, total)
	fmt.Println(StyleHighlight.Render(prefix) + " " + StyleValue.Render(msg))
}
