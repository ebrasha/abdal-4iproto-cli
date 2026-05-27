/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : countdown.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Pre-run countdown shown in argument mode allowing the
 *                operator to cancel the workflow by pressing 'q'.
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
	"time"

	"github.com/eiannone/keyboard"
)

// RunCountdown shows a countdown for `seconds` seconds. If the user presses
// 'q' (or Ctrl+C) during this window the function returns true (cancelled),
// otherwise it returns false at the end of the countdown.
func RunCountdown(seconds int) bool {
	if seconds <= 0 {
		return false
	}

	if err := keyboard.Open(); err != nil {
		// If raw keyboard access is unavailable we just sleep the same amount
		// of time so the user still sees the banner without an early abort.
		fmt.Println(StyleWarning.Render("[!] Cannot capture single key (q to cancel) – continuing after delay."))
		time.Sleep(time.Duration(seconds) * time.Second)
		return false
	}
	defer keyboard.Close()

	cancelled := make(chan bool, 1)
	done := make(chan struct{})

	// Watch keystrokes in the background.
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				char, key, err := keyboard.GetKey()
				if err != nil {
					return
				}
				if char == 'q' || char == 'Q' || key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
					cancelled <- true
					return
				}
			}
		}
	}()

	for remaining := seconds; remaining > 0; remaining-- {
		line := fmt.Sprintf("\r%s %s  %s",
			StyleTitle.Render(">>"),
			StyleHighlight.Render(fmt.Sprintf("Starting in %d second(s)...", remaining)),
			StyleHint.Render("[press 'q' to cancel]"),
		)
		fmt.Print(line + "   ")

		select {
		case <-cancelled:
			fmt.Println()
			fmt.Println(StyleWarning.Render("[!] Cancelled by user. Exiting safely."))
			close(done)
			return true
		case <-time.After(1 * time.Second):
		}
	}
	close(done)
	fmt.Println()
	fmt.Println(StyleSuccess.Render("[+] Starting now."))
	fmt.Println()
	return false
}
