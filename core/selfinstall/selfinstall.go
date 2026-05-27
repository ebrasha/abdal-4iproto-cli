/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : selfinstall.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Installs the abdal-4iproto-cli binary into a system-wide
 *                location so it can be invoked as a single command.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package selfinstall

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/ui"
)

// Install copies the running executable to a directory on PATH.
func Install() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve current executable: %w", err)
	}

	destDir, destName, err := targetPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("create destination directory: %w", err)
	}

	dest := filepath.Join(destDir, destName)
	if err := copyFile(exe, dest); err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(dest, 0o755); err != nil {
			return fmt.Errorf("chmod destination: %w", err)
		}
	}

	ui.SuccessBox("CLI Self-Install", fmt.Sprintf("Installed '%s' to:\n%s", config.AppCommandName, dest))
	return nil
}

func targetPath() (dir string, name string, err error) {
	name = config.AppCommandName
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	switch runtime.GOOS {
	case "linux":
		return "/usr/local/bin", name, nil
	case "windows":
		// Prefer a machine-wide location; fall back to user-local bin.
		programFiles := os.Getenv("ProgramFiles")
		if programFiles != "" {
			return filepath.Join(programFiles, "Abdal", "4iProto", "Cli"), name, nil
		}
		local := os.Getenv("LOCALAPPDATA")
		if local == "" {
			return "", "", fmt.Errorf("cannot resolve install directory on Windows")
		}
		return filepath.Join(local, "Programs", "Abdal", "4iProto", "Cli"), name, nil
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
