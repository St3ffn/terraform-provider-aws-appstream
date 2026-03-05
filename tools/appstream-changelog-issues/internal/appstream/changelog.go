// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package appstream

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"
)

type unused struct{}

const (
	// DefaultChangelogURL points to the upstream AWS SDK for Go v2 AppStream changelog.
	DefaultChangelogURL = "https://raw.githubusercontent.com/aws/aws-sdk-go-v2/main/service/appstream/CHANGELOG.md"
)

var (
	headingRE = regexp.MustCompile(`^# (v\d+\.\d+\.\d+) \(([^)]+)\)$`)

	allowedChangelogHosts = map[string]struct{}{
		"raw.githubusercontent.com": {},
	}
)

// Release represents one version section in the AppStream changelog.
type Release struct {
	Version string
	Date    string
	Notes   []string
}

// FetchChangelog downloads the AppStream changelog markdown from the provided URL.
// If client is nil, http.DefaultClient is used.
func FetchChangelog(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	if strings.TrimSpace(url) == "" {
		url = DefaultChangelogURL
	}

	url, err := sanitizeChangelogURL(url)
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// #nosec G704 -- request target is constrained by sanitizeChangelogURL to https://raw.githubusercontent.com only.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch changelog: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return body, nil
}

// ParseChangelog parses the AppStream changelog markdown into releases in the same
// order as present in the file (typically newest first).
func ParseChangelog(markdown []byte) ([]Release, error) {
	lines := strings.Split(string(markdown), "\n")
	releases := make([]Release, 0, 16)

	var current *Release
	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")

		if matches := headingRE.FindStringSubmatch(line); len(matches) == 3 {
			if current != nil {
				releases = append(releases, *current)
			}
			current = &Release{
				Version: matches[1],
				Date:    matches[2],
				Notes:   make([]string, 0, 8),
			}
			continue
		}

		if current == nil {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "* ") {
			current.Notes = append(current.Notes, strings.TrimSpace(strings.TrimPrefix(trimmed, "* ")))
			continue
		}

		// Support wrapped bullet lines.
		if len(current.Notes) > 0 && (strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t")) {
			last := len(current.Notes) - 1
			current.Notes[last] = current.Notes[last] + " " + trimmed
		}
	}

	if current != nil {
		releases = append(releases, *current)
	}

	if len(releases) == 0 {
		return nil, errors.New("no release sections found")
	}

	return releases, nil
}

// NewerThan returns releases with versions newer than the provided currentVersion.
func NewerThan(releases []Release, currentVersion string) ([]Release, error) {
	if !semver.IsValid(currentVersion) {
		return nil, fmt.Errorf("invalid current version %q", currentVersion)
	}

	filtered := make([]Release, 0, len(releases))
	for _, release := range releases {
		if !semver.IsValid(release.Version) {
			return nil, fmt.Errorf("invalid release version %q", release.Version)
		}
		if semver.Compare(release.Version, currentVersion) > 0 {
			filtered = append(filtered, release)
		}
	}

	return filtered, nil
}

// FilterFeatureReleases returns only releases that contain at least one
// changelog note marked as "**Feature**:".
func FilterFeatureReleases(releases []Release) []Release {
	filtered := make([]Release, 0, len(releases))

	for _, release := range releases {
		for _, note := range release.Notes {
			n := strings.TrimSpace(note)
			if strings.HasPrefix(strings.ToLower(n), strings.ToLower("**Feature**:")) {
				filtered = append(filtered, release)
				break
			}
		}
	}

	return filtered
}

func sanitizeChangelogURL(rawURL string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", fmt.Errorf("invalid changelog URL %q: %w", rawURL, err)
	}

	if !u.IsAbs() {
		return "", fmt.Errorf("invalid changelog URL %q: URL must be absolute", rawURL)
	}

	if !strings.EqualFold(u.Scheme, "https") {
		return "", fmt.Errorf("invalid changelog URL %q: only https URLs are allowed", rawURL)
	}

	if u.User != nil {
		return "", fmt.Errorf("invalid changelog URL %q: userinfo is not allowed", rawURL)
	}

	if u.Port() != "" {
		return "", fmt.Errorf("invalid changelog URL %q: custom ports are not allowed", rawURL)
	}

	host := strings.ToLower(u.Hostname())
	if _, ok := allowedChangelogHosts[host]; !ok {
		return "", fmt.Errorf("invalid changelog URL %q: host %q is not allowed", rawURL, host)
	}

	return u.String(), nil
}
