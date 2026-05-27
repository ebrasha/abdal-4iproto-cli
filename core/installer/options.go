/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : options.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Shared install option structs consumed by both the
 *                interactive workflow and Cobra flag-based commands.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package installer

import "abdal-4iproto-cli/core/keygen"

// Target describes which components to install.
type Target int

const (
	TargetAll Target = iota
	TargetServer
	TargetPanel
)

// Options bundles every user-tunable parameter for an installation run.
type Options struct {
	Target          Target
	ServerPorts     []int
	PanelPort       int
	PanelUsername   string
	PanelPassword   string
	Keygen          keygen.Options
	InstallServices bool
	SkipKeygen      bool
}

// DefaultOptions returns sane defaults for a full stack install.
func DefaultOptions() Options {
	return Options{
		Target:          TargetAll,
		PanelPort:       52202,
		PanelUsername:   "ebrasha",
		PanelPassword:   "ebrasha1366",
		Keygen:          keygen.DefaultOptions(),
		InstallServices: true,
	}
}
