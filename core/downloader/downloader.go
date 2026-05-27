/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : downloader.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Streams release assets from GitHub to disk with a live
 *                progress bar and verifies file integrity via SHA-256.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/github"
	"abdal-4iproto-cli/core/ui"
)

// Result describes a finished download operation.
type Result struct {
	FinalPath string
	Bytes     int64
	SHA256    string
}

// DownloadAsset streams the given asset to destDir using its canonical
// FinalLocal name. After download it verifies the SHA-256 digest if the
// release exposed one.
func DownloadAsset(sel *github.SelectedAsset, destDir string) (*Result, error) {
	if sel == nil || sel.Asset.DownloadURL == "" {
		return nil, fmt.Errorf("invalid asset selection")
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, fmt.Errorf("create destination directory: %w", err)
	}

	finalPath := filepath.Join(destDir, sel.FinalLocal)
	tmpPath := finalPath + ".part"

	client := &http.Client{Timeout: time.Duration(config.DownloadTimeoutSeconds) * time.Second}
	req, err := http.NewRequest(http.MethodGet, sel.Asset.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build download request: %w", err)
	}
	req.Header.Set("User-Agent", config.UserAgentHeader)
	req.Header.Set("Accept", config.GithubAcceptHeaderBinary)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("download HTTP %d for asset '%s'", resp.StatusCode, sel.Asset.Name)
	}

	size := sel.Asset.Size
	if size == 0 {
		size = resp.ContentLength
	}

	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open destination file: %w", err)
	}

	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription(ui.StyleHighlight.Render("Downloading ")+ui.StyleValue.Render(sel.Asset.Name)),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]▓[reset]",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionThrottle(80*time.Millisecond),
		progressbar.OptionClearOnFinish(),
	)

	hasher := sha256.New()
	mw := io.MultiWriter(out, hasher, bar)

	if _, err := io.Copy(mw, resp.Body); err != nil {
		out.Close()
		os.Remove(tmpPath)
		return nil, fmt.Errorf("stream body: %w", err)
	}
	if err := out.Close(); err != nil {
		os.Remove(tmpPath)
		return nil, fmt.Errorf("close output file: %w", err)
	}
	_ = bar.Finish()

	hash := hex.EncodeToString(hasher.Sum(nil))

	// Verify checksum when the API exposes one.
	if sel.SHA256 != "" {
		if !strings.EqualFold(hash, sel.SHA256) {
			os.Remove(tmpPath)
			return nil, fmt.Errorf("checksum mismatch for '%s': expected %s, got %s", sel.Asset.Name, sel.SHA256, hash)
		}
		ui.Success(fmt.Sprintf("SHA-256 verified for %s", sel.Asset.Name))
	} else {
		ui.Warning(fmt.Sprintf("Release did not expose a digest for %s – skipped checksum check.", sel.Asset.Name))
	}

	// Replace destination atomically.
	_ = os.Remove(finalPath)
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return nil, fmt.Errorf("rename temp file: %w", err)
	}

	// Ensure executable bit on Linux.
	if !sel.IsWindows {
		if err := os.Chmod(finalPath, 0o755); err != nil {
			ui.Warning(fmt.Sprintf("Could not chmod 0755 on %s: %v", finalPath, err))
		}
	}

	stat, _ := os.Stat(finalPath)
	var bytes int64
	if stat != nil {
		bytes = stat.Size()
	}
	return &Result{FinalPath: finalPath, Bytes: bytes, SHA256: hash}, nil
}
