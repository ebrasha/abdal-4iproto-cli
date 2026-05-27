/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : update.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 23:00:00
 * Description  : Renders the neon-green "new version available" notice
 *                using the same colour as the ASCII banner.
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

// UpdateNotice prints a green-bordered box advertising the new version
// and pointing the user to the release page.
func UpdateNotice(current, latest, releaseURL string) {
	title := "Update Available"
	body := fmt.Sprintf(
		"A new version of Abdal 4iProto Cli is available.\nCurrent : %s\nLatest  : %s\nUpdate at: %s",
		current, latest, releaseURL,
	)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorGreen).
		Foreground(ColorGreen).
		Bold(true).
		Padding(0, 2).
		MarginTop(1).
		MarginBottom(1)

	header := lipgloss.NewStyle().Foreground(ColorGreen).Bold(true).Render("⬆ " + title)
	separator := strings.Repeat("─", lipgloss.Width(title)+2)
	content := header + "\n" + separator + "\n" + body
	fmt.Println(style.Render(content))
}
