/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : selfinstall.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Cobra command that installs the CLI binary as the
 *                global abdal-4iproto-cli command.
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

	"abdal-4iproto-cli/core/selfinstall"
)

func newSelfInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "self-install",
		Short: "Install abdal-4iproto-cli into the system PATH",
		RunE: func(cmd *cobra.Command, args []string) error {
			return selfinstall.Install()
		},
	}
}
