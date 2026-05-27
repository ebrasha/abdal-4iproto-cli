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

// Install copies the running executable to a directory on PATH. If the
// destination file already exists it is overwritten; on Windows the
// existing binary may be locked by another process, so the function
// first moves it aside before writing the new copy.
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

	// Force-overwrite: drop any previous binary at this path.
	if err := prepareOverwrite(dest); err != nil {
		return err
	}

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

// prepareOverwrite makes sure a fresh copy can be written at dest by
// removing or moving aside any existing file. The renamed leftover
// (".old") is cleared on a best-effort basis so a Windows binary that
// is currently locked by another process never blocks the install.
func prepareOverwrite(dest string) error {
	info, err := os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("inspect destination: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("destination path is a directory, not a file: %s", dest)
	}

	if removeErr := os.Remove(dest); removeErr == nil {
		ui.Info("Existing binary overwritten: " + dest)
		return nil
	}

	// Likely Windows file-in-use: move the old binary aside and continue.
	stale := dest + ".old"
	_ = os.Remove(stale)
	if renameErr := os.Rename(dest, stale); renameErr != nil {
		return fmt.Errorf("cannot overwrite existing binary at %s (it may be locked by another process): %w", dest, renameErr)
	}
	ui.Info("Existing binary moved aside: " + stale)
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
