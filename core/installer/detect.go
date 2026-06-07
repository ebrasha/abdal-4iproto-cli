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
// Kept for callers that need a coarse "anything left behind" answer.
func (r ExistingInstallReport) IsPresent() bool {
	return r.HasServerBin || r.HasPanelBin || r.HasKeygenBin || r.HasConfigs
}

// IsServerPresent reports whether the server-side artefacts exist on disk.
func (r ExistingInstallReport) IsServerPresent() bool {
	return r.HasServerBin
}

// IsPanelPresent reports whether the panel-side artefacts exist on disk.
func (r ExistingInstallReport) IsPanelPresent() bool {
	return r.HasPanelBin
}

// IsFullStackPresent returns true only when *every* binary of the full
// stack (server, panel and keygen) is already installed, so the operator
// is never blocked just because a single partial component was left over.
func (r ExistingInstallReport) IsFullStackPresent() bool {
	return r.HasServerBin && r.HasPanelBin && r.HasKeygenBin
}

// MatchesTarget tells whether the artefacts that belong to the requested
// install scope already exist. The interactive flow uses this to decide
// when the fresh-install confirmation prompt is appropriate.
func (r ExistingInstallReport) MatchesTarget(target Target) bool {
	switch target {
	case TargetServer:
		return r.IsServerPresent()
	case TargetPanel:
		return r.IsPanelPresent()
	default:
		return r.IsFullStackPresent()
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
	rep.HasConfigs = pathExists(filepath.Join(dir, config.ServerConfigFileName)) ||
		pathExists(filepath.Join(dir, config.PanelConfigFileName)) ||
		pathExists(filepath.Join(dir, config.UsersFileName))
	return rep, nil
}

// FreshWipe stops every registered service and deletes the install
// directory so a full-stack Force install starts from a clean slate.
// Use FreshWipeFor when the operator picked a partial scope.
func FreshWipe() error {
	ui.Info("Performing fresh-install wipe: stopping services and clearing the installation directory.")
	_ = service.Uninstall(service.ComponentPanel)
	_ = service.Uninstall(service.ComponentServer)
	return uninstaller.Run(uninstaller.TargetAll, true)
}

// FreshWipeFor performs a scope-aware fresh wipe. The full-stack scope
// drops the entire install directory, while partial scopes (server or
// panel) clean only their own files so the co-installed component keeps
// working without losing its data.
func FreshWipeFor(target Target) error {
	switch target {
	case TargetServer:
		return freshWipeServer()
	case TargetPanel:
		return freshWipePanel()
	default:
		return FreshWipe()
	}
}

// freshWipeServer removes the server service plus every server-side file
// (server and keygen binaries, generated SSH keys, server_config.json,
// users.json and blocked_ips.json). The shared install directory and any
// panel artefacts are left intact.
func freshWipeServer() error {
	ui.Info("Performing fresh-install wipe for the server scope: stopping the server service and clearing server files.")
	_ = service.Uninstall(service.ComponentServer)

	dir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	candidates := []string{
		paths.ServerBinaryPath(dir),
		paths.KeygenBinaryPath(dir),
		filepath.Join(dir, config.ServerConfigFileName),
		filepath.Join(dir, config.UsersFileName),
		filepath.Join(dir, config.BlockedIPsFileName),
		filepath.Join(dir, config.DefaultKeyBaseName),
		filepath.Join(dir, config.DefaultKeyPublicName),
	}
	removeIfPresent(candidates)
	ui.Success("Server artefacts cleared. Panel binary and configuration kept intact.")
	return nil
}

// freshWipePanel removes the panel service and its files (panel binary
// and abdal-4iproto-panel.json). The shared install directory and any
// server artefacts (binaries, keys, users.json) are left intact.
func freshWipePanel() error {
	ui.Info("Performing fresh-install wipe for the panel scope: stopping the panel service and clearing panel files.")
	_ = service.Uninstall(service.ComponentPanel)

	dir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	candidates := []string{
		paths.PanelBinaryPath(dir),
		filepath.Join(dir, config.PanelConfigFileName),
	}
	removeIfPresent(candidates)
	ui.Success("Panel artefacts cleared. Server binary, keys and configuration kept intact.")
	return nil
}

// BinaryRefreshFor prepares a binary-only reinstall for the requested
// scope. Unlike FreshWipeFor it never deletes configuration files, SSH
// keys or user accounts: it merely stops the relevant service(s) and
// removes the executables so the downloader can drop fresh copies in
// place (the stop step is essential on Windows, where a running service
// keeps the .exe locked).
func BinaryRefreshFor(target Target) error {
	switch target {
	case TargetServer:
		return binaryRefreshServer()
	case TargetPanel:
		return binaryRefreshPanel()
	default:
		if err := binaryRefreshServer(); err != nil {
			return err
		}
		return binaryRefreshPanel()
	}
}

// binaryRefreshServer stops the server service and removes only the
// server and keygen executables. Configuration, keys and users persist.
func binaryRefreshServer() error {
	ui.Info("Reinstalling server binaries only: stopping the server service and replacing executables (configuration, keys and users are preserved).")
	_ = service.Stop(service.ComponentServer)

	dir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	removeIfPresent([]string{
		paths.ServerBinaryPath(dir),
		paths.KeygenBinaryPath(dir),
	})
	return nil
}

// binaryRefreshPanel stops the panel service and removes only the panel
// executable. The panel configuration file is preserved.
func binaryRefreshPanel() error {
	ui.Info("Reinstalling panel binary only: stopping the panel service and replacing the executable (panel configuration is preserved).")
	_ = service.Stop(service.ComponentPanel)

	dir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	removeIfPresent([]string{
		paths.PanelBinaryPath(dir),
	})
	return nil
}

// removeIfPresent silently deletes every existing file in the list;
// missing entries are ignored because partial installs may not contain
// every artefact we know about.
func removeIfPresent(files []string) {
	for _, p := range files {
		if _, err := os.Stat(p); err == nil {
			_ = os.Remove(p)
		}
	}
}

func pathExists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}
