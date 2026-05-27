/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : paths.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Resolves installation directories and canonical binary
 *                names for the current operating system.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"abdal-4iproto-cli/core/config"
)

// InstallDir returns the absolute installation directory for the server stack.
func InstallDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		local := os.Getenv("LOCALAPPDATA")
		if local == "" {
			return "", os.ErrNotExist
		}
		return filepath.Join(local, config.WindowsInstallDirRelative), nil
	case "linux":
		return config.LinuxInstallDir, nil
	default:
		return "", os.ErrInvalid
	}
}

// ServerBinaryName returns the canonical server binary file name.
func ServerBinaryName() string {
	if runtime.GOOS == "windows" {
		return config.WindowsServerBinaryName + ".exe"
	}
	return config.LinuxServerBinaryName
}

// PanelBinaryName returns the canonical panel binary file name.
func PanelBinaryName() string {
	if runtime.GOOS == "windows" {
		return config.WindowsPanelBinaryName + ".exe"
	}
	return config.LinuxPanelBinaryName
}

// KeygenBinaryName returns the canonical keygen binary file name.
func KeygenBinaryName() string {
	if runtime.GOOS == "windows" {
		return config.WindowsKeygenBinaryName + ".exe"
	}
	return config.LinuxKeygenBinaryName
}

// ServerBinaryPath joins install dir with the server binary name.
func ServerBinaryPath(installDir string) string {
	return filepath.Join(installDir, ServerBinaryName())
}

// PanelBinaryPath joins install dir with the panel binary name.
func PanelBinaryPath(installDir string) string {
	return filepath.Join(installDir, PanelBinaryName())
}

// KeygenBinaryPath joins install dir with the keygen binary name.
func KeygenBinaryPath(installDir string) string {
	return filepath.Join(installDir, KeygenBinaryName())
}
