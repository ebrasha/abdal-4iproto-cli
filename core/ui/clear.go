/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : clear.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-30 18:00:00
 * Description  : Cross-platform screen clearing helpers used between
 *                interactive menu transitions.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ClearScreen wipes the terminal so each menu transition starts from a
// clean slate. On Windows the legacy console may not honour ANSI escape
// codes, so we shell out to "cmd /c cls"; everywhere else we send the
// standard VT100 sequence which also clears the scrollback buffer.
func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			// Fallback to ANSI in case cmd is unavailable.
			fmt.Print("\033[H\033[2J\033[3J")
		}
	default:
		fmt.Print("\033[H\033[2J\033[3J")
	}
}

// ClearAndBanner clears the screen and re-prints the program banner so
// the user always sees the identity of the tool at the top of each
// menu screen.
func ClearAndBanner() {
	ClearScreen()
	PrintBanner()
}

// PressEnter pauses the workflow until the operator hits Enter. Use it
// right before a screen clear so the previous command output is not
// wiped away faster than the user can read it.
func PressEnter() {
	fmt.Println()
	fmt.Println(StyleHint.Render("Press Enter to return to the menu..."))
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
