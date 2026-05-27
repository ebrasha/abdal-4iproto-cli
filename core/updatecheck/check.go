/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : check.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 23:00:00
 * Description  : Runs the self-update probe at start-up and renders the
 *                neon-green notice when a newer release is detected.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package updatecheck

import (
	gh "abdal-4iproto-cli/core/github"
	"abdal-4iproto-cli/core/ui"
)

// Notify queries the GitHub releases API for a newer CLI version and
// prints a colored notice when one is available. Failures are silenced
// so the rest of the workflow is never blocked by network problems.
func Notify() {
	info, err := gh.CheckCliUpdate()
	if err != nil || info == nil {
		return
	}
	if info.UpdateNeeded {
		ui.UpdateNotice(info.Current, info.Latest, info.ReleaseURL)
	}
}
