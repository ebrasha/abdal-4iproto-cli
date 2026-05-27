/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : update.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 23:00:00
 * Description  : Self-update checker that compares the running CLI
 *                version against the latest published GitHub release.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package github

import (
	"strconv"
	"strings"

	"abdal-4iproto-cli/core/config"
)

// UpdateInfo describes the outcome of a self-update check.
type UpdateInfo struct {
	Current      string
	Latest       string
	UpdateNeeded bool
	ReleaseURL   string
}

// CheckCliUpdate contacts the GitHub releases API for the CLI itself and
// returns whether the running binary is older than the latest published
// release. Network or parsing failures are returned unchanged so the
// caller can decide whether to display them.
func CheckCliUpdate() (*UpdateInfo, error) {
	rel, err := FetchLatestRelease(config.CliLatestReleaseAPI)
	if err != nil {
		return nil, err
	}
	latest := normalizeTag(rel.TagName)
	current := normalizeTag(config.AppVersion)
	return &UpdateInfo{
		Current:      current,
		Latest:       latest,
		UpdateNeeded: compareVersions(current, latest) < 0,
		ReleaseURL:   rel.HTMLURL,
	}, nil
}

// normalizeTag trims optional "v" prefix and surrounding whitespace so a
// tag like "v1.5" or " 1.5 " becomes "1.5".
func normalizeTag(tag string) string {
	tag = strings.TrimSpace(tag)
	tag = strings.TrimPrefix(tag, "v")
	tag = strings.TrimPrefix(tag, "V")
	return tag
}

// compareVersions returns -1, 0, or +1 when a is older, equal, or newer
// than b. Each version is split on '.' and compared numerically segment
// by segment so "1.5" < "1.6" and "1.5" == "1.5".
func compareVersions(a, b string) int {
	pa := splitVersion(a)
	pb := splitVersion(b)
	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}
	for i := 0; i < n; i++ {
		va := segmentAt(pa, i)
		vb := segmentAt(pb, i)
		if va < vb {
			return -1
		}
		if va > vb {
			return 1
		}
	}
	return 0
}

// splitVersion converts a dotted version string into integer segments,
// silently substituting zero for non-numeric parts so a stray suffix
// cannot crash the comparison.
func splitVersion(v string) []int {
	parts := strings.Split(v, ".")
	out := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			n = 0
		}
		out[i] = n
	}
	return out
}

// segmentAt returns the i-th segment or 0 if the version has fewer parts.
func segmentAt(segments []int, i int) int {
	if i >= len(segments) {
		return 0
	}
	return segments[i]
}
