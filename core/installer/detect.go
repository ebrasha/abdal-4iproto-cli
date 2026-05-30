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
// Configs are split into server-side and panel-side so the workflow can
// reason about each component independently.
type ExistingInstallReport struct {
	InstallDir       string
	HasServerBin     bool
	HasPanelBin      bool
	HasKeygenBin     bool
	HasServerConfigs bool
	HasPanelConfig   bool
	// HasConfigs preserves the original "any config file present" flag
	// for callers that do not care about the per-component split.
	HasConfigs bool
}

// IsPresent reports whether any installation artefact was found.
func (r ExistingInstallReport) IsPresent() bool {
	return r.HasServerBin || r.HasPanelBin || r.HasKeygenBin || r.HasConfigs
}

// HasServerStack reports whether any server-side artefact (binary,
// keygen helper, server_config.json, blocked_ips.json or users.json)
// is present on disk. This is the unit that "Server only" installs and
// uninstalls operate on.
func (r ExistingInstallReport) HasServerStack() bool {
	return r.HasServerBin || r.HasKeygenBin || r.HasServerConfigs
}

// HasPanelStack reports whether any panel-side artefact (binary or
// abdal-4iproto-panel.json) is present on disk.
func (r ExistingInstallReport) HasPanelStack() bool {
	return r.HasPanelBin || r.HasPanelConfig
}

// IsTargetPresent reports whether the components that belong to the
// requested install scope already exist on disk. It is used to decide
// per-scope whether a "fresh install" confirmation is required.
func (r ExistingInstallReport) IsTargetPresent(t Target) bool {
	switch t {
	case TargetServer:
		return r.HasServerStack()
	case TargetPanel:
		return r.HasPanelStack()
	default:
		return r.HasServerStack() && r.HasPanelStack()
	}
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
	rep.HasServerConfigs = pathExists(filepath.Join(dir, config.ServerConfigFileName)) ||
		pathExists(filepath.Join(dir, config.UsersFileName)) ||
		pathExists(filepath.Join(dir, config.BlockedIPsFileName))
	rep.HasPanelConfig = pathExists(filepath.Join(dir, config.PanelConfigFileName))
	rep.HasConfigs = rep.HasServerConfigs || rep.HasPanelConfig
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

// FreshWipeTarget removes only the artefacts that belong to the
// requested install scope so a fresh re-install of one component does
// not destroy the other. For TargetAll it delegates to FreshWipe.
func FreshWipeTarget(t Target) error {
	switch t {
	case TargetServer:
		return freshWipeServer()
	case TargetPanel:
		return freshWipePanel()
	default:
		return FreshWipe()
	}
}

// freshWipeServer removes the server service, server/keygen binaries
// and the server-side config files but keeps the panel artefacts intact.
func freshWipeServer() error {
	ui.Info("Performing server-only fresh wipe: stopping server service and clearing server files.")
	_ = service.Uninstall(service.ComponentServer)

	dir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Remove the server runtime, keygen helper and any generated keys.
	removeIfExists(paths.ServerBinaryPath(dir))
	removeIfExists(paths.KeygenBinaryPath(dir))
	removeIfExists(filepath.Join(dir, config.DefaultKeyBaseName))
	removeIfExists(filepath.Join(dir, config.DefaultKeyPublicName))

	// Wipe the JSON files that belong exclusively to the server.
	removeIfExists(filepath.Join(dir, config.ServerConfigFileName))
	removeIfExists(filepath.Join(dir, config.UsersFileName))
	removeIfExists(filepath.Join(dir, config.BlockedIPsFileName))
	return nil
}

// freshWipePanel removes the panel service, panel binary and the panel
// config file while leaving server-side files untouched.
func freshWipePanel() error {
	ui.Info("Performing panel-only fresh wipe: stopping panel service and clearing panel files.")
	_ = service.Uninstall(service.ComponentPanel)

	dir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	removeIfExists(paths.PanelBinaryPath(dir))
	removeIfExists(filepath.Join(dir, config.PanelConfigFileName))
	return nil
}

// removeIfExists deletes a file or directory when present and ignores
// the "not exist" case so the caller can stay declarative.
func removeIfExists(p string) {
	if p == "" {
		return
	}
	if _, err := os.Stat(p); err != nil {
		return
	}
	_ = os.RemoveAll(p)
}

func pathExists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}
