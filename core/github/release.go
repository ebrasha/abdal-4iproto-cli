/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : release.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Talks to the GitHub REST API to fetch the latest
 *                release metadata, picks the asset that matches the
 *                runtime OS/architecture, and exposes its SHA256 digest.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"abdal-4iproto-cli/core/config"
)

// Asset represents a single GitHub release asset as returned by the API.
type Asset struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"browser_download_url"`
	APIURL      string `json:"url"`
	ContentType string `json:"content_type"`
	Digest      string `json:"digest"` // Format: "sha256:<hex>"
}

// Release contains the subset of fields we care about from a GitHub release.
type Release struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	PublishedAt string  `json:"published_at"`
	Body        string  `json:"body"`
	HTMLURL     string  `json:"html_url"`
	Assets      []Asset `json:"assets"`
}

// SelectedAsset bundles the chosen asset and the parsed checksum.
type SelectedAsset struct {
	Asset      Asset
	SHA256     string // Lower-case hex, no prefix.
	OSName     string
	Arch       string
	IsWindows  bool
	FinalLocal string // Canonical local file name to use after rename.
}

// FetchLatestRelease retrieves the latest published release for the given
// repository endpoint.
func FetchLatestRelease(endpoint string) (*Release, error) {
	client := &http.Client{Timeout: time.Duration(config.HTTPTimeoutSeconds) * time.Second}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", config.GithubAcceptHeaderJSON)
	req.Header.Set("X-GitHub-Api-Version", config.GithubAPIVersion)
	req.Header.Set("User-Agent", config.UserAgentHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decode release JSON: %w", err)
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("no tag_name in latest release response")
	}
	return &rel, nil
}

// ChooseAsset picks the asset that matches runtime OS/arch by applying the
// provided pattern (containing the "{arch}" placeholder). When no exact match
// is found it falls back to a case-insensitive substring search.
func ChooseAsset(rel *Release, windowsPattern, linuxPattern, finalWindowsName, finalLinuxName string) (*SelectedAsset, error) {
	if rel == nil {
		return nil, fmt.Errorf("nil release")
	}

	osName := runtime.GOOS
	arch := runtime.GOARCH

	var pattern, finalName string
	switch osName {
	case "windows":
		pattern = strings.ReplaceAll(windowsPattern, "{arch}", arch)
		finalName = finalWindowsName + ".exe"
	case "linux":
		pattern = strings.ReplaceAll(linuxPattern, "{arch}", arch)
		finalName = finalLinuxName
	default:
		return nil, fmt.Errorf("unsupported OS: %s", osName)
	}

	// Pass 1: exact (case-insensitive) match on the pattern.
	for i := range rel.Assets {
		if strings.EqualFold(rel.Assets[i].Name, pattern) {
			return buildSelected(rel.Assets[i], osName, arch, finalName)
		}
	}

	// Pass 2: try matching with the arch token present anywhere along with the
	// OS keyword (helps with releases that add suffixes).
	osKey := osName
	for i := range rel.Assets {
		lname := strings.ToLower(rel.Assets[i].Name)
		if strings.Contains(lname, osKey) && strings.Contains(lname, arch) {
			return buildSelected(rel.Assets[i], osName, arch, finalName)
		}
	}

	return nil, fmt.Errorf("no release asset found matching '%s' (or '%s' + '%s') for tag %s", pattern, osKey, arch, rel.TagName)
}

// buildSelected normalizes the asset's digest field and constructs a
// SelectedAsset.
func buildSelected(a Asset, osName, arch, finalName string) (*SelectedAsset, error) {
	digest := strings.TrimSpace(a.Digest)
	var sha string
	if digest != "" {
		// The API returns "sha256:abcdef..." per GitHub documentation.
		if strings.Contains(digest, ":") {
			parts := strings.SplitN(digest, ":", 2)
			algo := strings.ToLower(strings.TrimSpace(parts[0]))
			value := strings.ToLower(strings.TrimSpace(parts[1]))
			if algo != "sha256" {
				return nil, fmt.Errorf("unsupported digest algorithm '%s' for asset '%s'", algo, a.Name)
			}
			sha = value
		} else {
			sha = strings.ToLower(digest)
		}
	}

	return &SelectedAsset{
		Asset:      a,
		SHA256:     sha,
		OSName:     osName,
		Arch:       arch,
		IsWindows:  osName == "windows",
		FinalLocal: finalName,
	}, nil
}
