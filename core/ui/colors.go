/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : colors.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Neon color palette and shared lipgloss styles used by
 *                every CLI surface (titles, banners, prompts, boxes).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package ui

import "github.com/charmbracelet/lipgloss"

// Neon color palette (xterm 256-color codes).
var (
	ColorRed     = lipgloss.Color("196") // Neon Red
	ColorGreen   = lipgloss.Color("46")  // Neon Green
	ColorYellow  = lipgloss.Color("226") // Neon Yellow
	ColorBlue    = lipgloss.Color("51")  // Neon Blue / Cyan
	ColorPurple  = lipgloss.Color("129") // Neon Purple
	ColorPink    = lipgloss.Color("201") // Neon Pink
	ColorOrange  = lipgloss.Color("208") // Neon Orange
	ColorWhite   = lipgloss.Color("15")  // Bright White
	ColorCyan    = lipgloss.Color("87")  // Bright Cyan
	ColorMagenta = lipgloss.Color("165") // Bright Magenta
	ColorGray    = lipgloss.Color("245") // Muted gray for hints
)

// StyleBanner is reserved for the ASCII art block.
var StyleBanner = lipgloss.NewStyle().Foreground(ColorGreen).Bold(true)

// StyleProgrammer is the credit line under the banner.
var StyleProgrammer = lipgloss.NewStyle().Foreground(ColorCyan).Bold(true)

// StyleTitle is used for section titles.
var StyleTitle = lipgloss.NewStyle().Foreground(ColorMagenta).Bold(true)

// StyleSuccess marks successful messages.
var StyleSuccess = lipgloss.NewStyle().Foreground(ColorGreen).Bold(true)

// StyleError marks fatal errors.
var StyleError = lipgloss.NewStyle().Foreground(ColorRed).Bold(true)

// StyleWarning marks recoverable warnings.
var StyleWarning = lipgloss.NewStyle().Foreground(ColorYellow).Bold(true)

// StyleInfo marks informational lines.
var StyleInfo = lipgloss.NewStyle().Foreground(ColorBlue)

// StyleHint formats subtle hints below prompts.
var StyleHint = lipgloss.NewStyle().Foreground(ColorGray).Italic(true)

// StyleHighlight emphasizes a critical value inside a sentence.
var StyleHighlight = lipgloss.NewStyle().Foreground(ColorOrange).Bold(true)

// StyleLabel formats keys in key/value summaries.
var StyleLabel = lipgloss.NewStyle().Foreground(ColorPurple).Bold(true)

// StyleValue formats the value part of key/value summaries.
var StyleValue = lipgloss.NewStyle().Foreground(ColorWhite)

// StylePrompt formats prompt question prefixes.
var StylePrompt = lipgloss.NewStyle().Foreground(ColorPink).Bold(true)
