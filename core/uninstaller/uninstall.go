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

// Run removes services and optionally deletes the install directory.
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
	case TargetPanel:
		if err := service.Uninstall(service.ComponentPanel); err != nil {
			return err
		}
	default:
		_ = service.Uninstall(service.ComponentPanel)
		_ = service.Uninstall(service.ComponentServer)
	}

	if removeFiles {
		if err := os.RemoveAll(installDir); err != nil {
			return fmt.Errorf("remove install directory: %w", err)
		}
		ui.Success("Removed installation directory: " + installDir)
	}

	ui.SuccessBox("Uninstall Complete", "Selected components have been removed.")
	return nil
}
