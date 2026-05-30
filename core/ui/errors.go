/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : errors.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-30 19:00:00
 * Description  : Shared sentinel errors used by interactive menus to
 *                signal navigation intents (e.g. Back) so parent
 *                handlers can skip post-action prompts cleanly.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package ui

import "errors"

// ErrUserBack signals that the operator chose to navigate back to the
// previous menu. Parent loops should treat this as a successful return
// and skip the post-action "Press Enter" pause, because the user has
// already decided to leave the current screen.
var ErrUserBack = errors.New("user navigated back")

// IsBack reports whether err represents a user-initiated back navigation.
func IsBack(err error) bool {
	return errors.Is(err, ErrUserBack)
}
