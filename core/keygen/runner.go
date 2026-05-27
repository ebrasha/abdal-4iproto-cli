/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : runner.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Executes the downloaded SSH key generator binary with
 *                configurable key type, size, force overwrite, and output.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package keygen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/filesgen"
	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/ui"
)

// Options describes the key generation parameters.
type Options struct {
	Type       string // rsa, ed25519, ecdsa
	Bits       int
	Force      bool
	OutputFile string // private key filename (public becomes <f>.pub)
}

// DefaultOptions returns the installer defaults (-force -t ed25519 -b 4096).
func DefaultOptions() Options {
	return Options{
		Type:       config.DefaultKeygenType,
		Bits:       config.DefaultKeygenBits,
		Force:      config.DefaultKeygenForce,
		OutputFile: config.DefaultKeyBaseName,
	}
}

// Generate runs the keygen binary inside installDir and returns the key names.
func Generate(installDir string, opts Options) (filesgen.KeyFileNames, error) {
	bin := paths.KeygenBinaryPath(installDir)
	if _, err := os.Stat(bin); err != nil {
		return filesgen.KeyFileNames{}, fmt.Errorf("keygen binary not found at %s: %w", bin, err)
	}

	if strings.TrimSpace(opts.OutputFile) == "" {
		opts.OutputFile = config.DefaultKeyBaseName
	}

	args := []string{
		"-t", opts.Type,
		"-b", fmt.Sprintf("%d", opts.Bits),
		"-f", opts.OutputFile,
	}
	if opts.Force {
		args = append(args, "-force")
	}

	ui.Info(fmt.Sprintf("Generating SSH keys (%s, %d bits) in %s", opts.Type, opts.Bits, installDir))

	cmd := exec.Command(bin, args...)
	cmd.Dir = installDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return filesgen.KeyFileNames{}, fmt.Errorf("keygen failed: %w", err)
	}

	private := opts.OutputFile
	public := private + ".pub"
	if _, err := os.Stat(filepath.Join(installDir, private)); err != nil {
		return filesgen.KeyFileNames{}, fmt.Errorf("private key not created: %w", err)
	}
	if _, err := os.Stat(filepath.Join(installDir, public)); err != nil {
		return filesgen.KeyFileNames{}, fmt.Errorf("public key not created: %w", err)
	}

	ui.Success(fmt.Sprintf("SSH keys created: %s / %s", private, public))
	return filesgen.KeyFileNames{Private: private, Public: public}, nil
}
