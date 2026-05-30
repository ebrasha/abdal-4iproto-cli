/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : profile.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-30 18:00:00
 * Description  : Forces a colored terminal profile so the neon banner
 *                renders identically on Windows, Linux, and macOS.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package ui

import "os"

// init ensures lipgloss/termenv do not silently downgrade to a no-color
// profile on Linux terminals that fail TTY/COLORTERM detection (some
// SSH sessions, minimal containers, CI runners, etc.). The hints below
// are honoured by termenv – the colour backend used by lipgloss – so
// the ASCII banner and bordered boxes stay vivid everywhere.
func init() {
	if os.Getenv("CLICOLOR_FORCE") == "" {
		_ = os.Setenv("CLICOLOR_FORCE", "1")
	}
	if os.Getenv("COLORTERM") == "" {
		_ = os.Setenv("COLORTERM", "truecolor")
	}
	if os.Getenv("TERM") == "" || os.Getenv("TERM") == "dumb" {
		_ = os.Setenv("TERM", "xterm-256color")
	}
}
