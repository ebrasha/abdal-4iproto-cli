/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : selfupdate.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-31 01:35:00
 * Description  : Cobra command that downloads the latest CLI release
 *                from GitHub, verifies its SHA-256 digest, and replaces
 *                the running binary in place.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package cli

import (
	"github.com/spf13/cobra"

	"abdal-4iproto-cli/core/selfupdate"
)

// newSelfUpdateCmd wires up the `abdal-4iproto-cli self-update` command
// which mirrors the "Update CLI Command" entry inside the interactive
// "Manage CLI Command" submenu.
func newSelfUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "self-update",
		Short: "Download and apply the latest abdal-4iproto-cli release from GitHub",
		RunE: func(cmd *cobra.Command, args []string) error {
			return selfupdate.Update()
		},
	}
}
