/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : uninstall.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Removes Abdal 4iProto services and installation files
 *                from the host operating system.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package uninstaller

import (
	"fmt"
	"os"

	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/service"
	"abdal-4iproto-cli/core/ui"
)

// Target mirrors installer.Target for uninstall scope.
type Target int

const (
	TargetAll Target = iota
	TargetServer
	TargetPanel
)

// Run removes services and, for a full-stack uninstall only, optionally
// deletes the entire installation directory. Partial uninstalls (server
// or panel only) never remove the shared install path so the remaining
// component keeps working.
func Run(target Target, removeFiles bool) error {
	installDir, err := paths.InstallDir()
	if err != nil {
		return fmt.Errorf("resolve install directory: %w", err)
	}

	ui.SectionHeader("Abdal 4iProto Uninstall")

	switch target {
	case TargetServer:
		if err := service.Uninstall(service.ComponentServer); err != nil {
			return err
		}
		ui.Info("Installation directory preserved: " + installDir)
	case TargetPanel:
		if err := service.Uninstall(service.ComponentPanel); err != nil {
			return err
		}
		ui.Info("Installation directory preserved: " + installDir)
	default:
		_ = service.Uninstall(service.ComponentPanel)
		_ = service.Uninstall(service.ComponentServer)

		// Only full-stack uninstall may delete the shared install folder.
		if removeFiles {
			if err := os.RemoveAll(installDir); err != nil {
				return fmt.Errorf("remove install directory: %w", err)
			}
			ui.Success("Removed installation directory: " + installDir)
		} else {
			ui.Info("Services removed; installation directory kept: " + installDir)
		}
	}

	ui.SuccessBox("Uninstall Complete", "Selected components have been removed.")
	return nil
}
