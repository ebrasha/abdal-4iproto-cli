/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : uninstall.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Cobra uninstall subcommand for full or partial removal.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"abdal-4iproto-cli/core/uninstaller"
)

func newUninstallCmd() *cobra.Command {
	var (
		serverOnly bool
		panelOnly  bool
		keepFiles  bool
	)

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove Abdal 4iProto components and services",
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverOnly && panelOnly {
				return fmt.Errorf("cannot use --server-only and --panel-only together")
			}
			var target uninstaller.Target
			switch {
			case serverOnly:
				target = uninstaller.TargetServer
			case panelOnly:
				target = uninstaller.TargetPanel
			default:
				target = uninstaller.TargetAll
			}
			// Delete install directory only on full-stack uninstall (not --keep-files).
			removeInstallDir := target == uninstaller.TargetAll && !keepFiles
			return uninstaller.Run(target, removeInstallDir)
		},
	}

	cmd.Flags().BoolVar(&serverOnly, "server-only", false, "Uninstall server only")
	cmd.Flags().BoolVar(&panelOnly, "panel-only", false, "Uninstall panel only")
	cmd.Flags().BoolVar(&keepFiles, "keep-files", false, "Remove services but keep installation files")

	return cmd
}
