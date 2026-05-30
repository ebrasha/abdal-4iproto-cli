/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : detect.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-30 18:00:00
 * Description  : Detects whether Abdal 4iProto is already installed so
 *                we can warn the operator before overwriting state.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package installer

import (
	"errors"
	"os"
	"path/filepath"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/service"
	"abdal-4iproto-cli/core/ui"
	"abdal-4iproto-cli/core/uninstaller"
)

// ErrAlreadyInstalled is returned by the install pipeline when an
// existing installation is detected and Force is false. The interactive
// menu turns this into a confirmation prompt; the CLI surfaces it as a
// hard error so scripts do not overwrite state silently.
var ErrAlreadyInstalled = errors.New("Abdal 4iProto is already installed on this host")

// ExistingInstallReport summarises the state of a previous installation.
type ExistingInstallReport struct {
	InstallDir   string
	HasServerBin bool
	HasPanelBin  bool
	HasKeygenBin bool
	HasConfigs   bool
}

// IsPresent reports whether any installation artefact was found.
func (r ExistingInstallReport) IsPresent() bool {
	return r.HasServerBin || r.HasPanelBin || r.HasKeygenBin || r.HasConfigs
}

// DetectExisting inspects the install directory for any leftover files
// from a previous run. Missing directory => empty report (not present).
func DetectExisting() (ExistingInstallReport, error) {
	dir, err := paths.InstallDir()
	if err != nil {
		return ExistingInstallReport{}, err
	}
	rep := ExistingInstallReport{InstallDir: dir}

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return rep, nil
		}
		return rep, err
	}

	rep.HasServerBin = pathExists(paths.ServerBinaryPath(dir))
	rep.HasPanelBin = pathExists(paths.PanelBinaryPath(dir))
	rep.HasKeygenBin = pathExists(paths.KeygenBinaryPath(dir))
	rep.HasConfigs = pathExists(filepath.Join(dir, config.ServerConfigFileName)) ||
		pathExists(filepath.Join(dir, config.PanelConfigFileName)) ||
		pathExists(filepath.Join(dir, config.UsersFileName))
	return rep, nil
}

// FreshWipe stops registered services and deletes the install directory
// so a Force install can start from a clean slate.
func FreshWipe() error {
	ui.Info("Performing fresh-install wipe: stopping services and clearing the installation directory.")
	_ = service.Uninstall(service.ComponentPanel)
	_ = service.Uninstall(service.ComponentServer)
	return uninstaller.Run(uninstaller.TargetAll, true)
}

func pathExists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}
