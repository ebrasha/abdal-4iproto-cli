/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : selfupdate.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-31 01:30:00
 * Description  : Self-update routine for the abdal-4iproto-cli binary.
 *                Downloads the matching GitHub release asset, verifies
 *                its SHA-256 digest, and swaps the running executable
 *                in place using the inconshreveable/go-update package.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package selfupdate

import (
	"crypto"
	_ "crypto/sha256" // register the SHA-256 implementation for go-update
	"encoding/hex"
	"fmt"
	"os"

	gu "github.com/inconshreveable/go-update"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/downloader"
	gh "abdal-4iproto-cli/core/github"
	"abdal-4iproto-cli/core/ui"
)

// Update probes the GitHub releases API for a newer CLI build and,
// when one is available, downloads the matching OS/arch asset, verifies
// its SHA-256 digest against the release metadata, and atomically
// replaces the currently running executable with the new copy.
//
// The currently executing binary stays usable until the very last step
// (Apply), so an aborted update never leaves the user without a CLI.
func Update() error {
	info, err := gh.CheckCliUpdate()
	if err != nil {
		return fmt.Errorf("check for updates: %w", err)
	}
	if info == nil {
		return fmt.Errorf("update check returned no data")
	}

	if !info.UpdateNeeded {
		ui.SuccessBox("Already Up-To-Date", fmt.Sprintf(
			"You are running the latest version of %s.\nCurrent : %s\nLatest  : %s",
			config.AppName, info.Current, info.Latest,
		))
		return nil
	}

	ui.SectionHeader("CLI Self-Update")
	ui.KeyValueBox("Update Plan", [][2]string{
		{"Application", config.AppName},
		{"Current Version", info.Current},
		{"Latest Version", info.Latest},
		{"Release Page", info.ReleaseURL},
	})

	// Step 1 – fetch the release metadata so we can pick the right asset.
	rel, err := gh.FetchLatestRelease(config.CliLatestReleaseAPI)
	if err != nil {
		return fmt.Errorf("fetch latest release: %w", err)
	}

	// Step 2 – select the asset that matches runtime.GOOS / GOARCH using
	// the documented naming patterns:
	//   Windows : abdal-4iproto-cli-windows-<arch>.exe
	//   Linux   : abdal_4iproto_cli_linux_<arch>
	sel, err := gh.ChooseAsset(rel,
		config.WindowsCliAssetPattern, config.LinuxCliAssetPattern,
		config.WindowsCliBinaryName, config.LinuxCliBinaryName,
	)
	if err != nil {
		return fmt.Errorf("select matching asset for this OS/arch: %w", err)
	}

	// Step 3 – download to a private temporary directory; the running
	// binary keeps working until the final atomic swap succeeds.
	tmpDir, err := os.MkdirTemp("", "abdal-4iproto-cli-update-*")
	if err != nil {
		return fmt.Errorf("create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.Step(1, 3, "Downloading new CLI binary")
	result, err := downloader.DownloadAsset(sel, tmpDir)
	if err != nil {
		return fmt.Errorf("download new binary: %w", err)
	}

	// Step 4 – the downloader already verified the SHA-256 against the
	// digest exposed by the GitHub API, but we surface and double-check
	// it here so the operator sees the integrity guarantee explicitly.
	ui.Step(2, 3, "Verifying SHA-256 checksum")
	if sel.SHA256 == "" {
		ui.Warning("Release did not advertise a SHA-256 digest – relying on the downloader's hash trace only.")
	} else {
		ui.Info(fmt.Sprintf("Expected SHA-256: %s", sel.SHA256))
	}
	ui.Info(fmt.Sprintf("Observed SHA-256: %s", result.SHA256))

	// Step 5 – swap the running executable with the freshly downloaded
	// file. go-update writes a sibling temporary file first so any
	// failure is rolled back automatically.
	ui.Step(3, 3, "Replacing the current binary")
	if err := applyBinary(result.FinalPath, sel.SHA256); err != nil {
		return fmt.Errorf("replace binary: %w", err)
	}

	exe, _ := os.Executable()
	ui.SuccessBox("Update Complete", fmt.Sprintf(
		"%s has been updated to version %s.\nBinary: %s\nRestart the command to use the new version.",
		config.AppName, info.Latest, exe,
	))
	return nil
}

// applyBinary performs the atomic swap. When the release exposes a
// SHA-256 digest we pass it to go-update so the package double-checks
// the file one more time before writing it over the running executable.
func applyBinary(downloadedPath, expectedSHA256 string) error {
	f, err := os.Open(downloadedPath)
	if err != nil {
		return fmt.Errorf("open downloaded file: %w", err)
	}
	defer f.Close()

	opts := gu.Options{}
	if expectedSHA256 != "" {
		checksum, decodeErr := hex.DecodeString(expectedSHA256)
		if decodeErr != nil {
			return fmt.Errorf("decode expected SHA-256: %w", decodeErr)
		}
		opts.Checksum = checksum
		opts.Hash = crypto.SHA256
	}

	if err := gu.Apply(f, opts); err != nil {
		if rerr := gu.RollbackError(err); rerr != nil {
			return fmt.Errorf("apply failed (%v); rollback also failed: %v", err, rerr)
		}
		return err
	}
	return nil
}
