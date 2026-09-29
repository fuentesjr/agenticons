// Package models finds newer point releases of pinned Codex models using the
// model catalog reported by the installed Codex CLI.
//
// A model family is its major version plus its variant suffix, so gpt-6-sol
// and gpt-6.1-sol share a family while gpt-7-sol and gpt-6-astra do not. Only
// point releases within a family count as upgrades; major versions and
// variant changes stay deliberate decisions.
package models

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Model is the subset of a Codex catalog entry needed to pick an upgrade.
type Model struct {
	Slug            string `json:"slug"`
	Visibility      string `json:"visibility"`
	ReasoningLevels []struct {
		Effort string `json:"effort"`
	} `json:"supported_reasoning_levels"`
}

// Catalog is the parsed Codex model cache.
type Catalog struct {
	Models []Model `json:"models"`
}

// slugRE splits a slug into major version, optional dotted point version, and
// variant suffix: gpt-6.1-sol -> "6", "1", "-sol".
var slugRE = regexp.MustCompile(`^gpt-(\d+)(?:\.(\d+(?:\.\d+)*))?(-[a-z0-9-]+)?$`)

// FetchCatalog asks the installed Codex CLI for its current model catalog.
// The CLI's on-disk cache is not used because any Codex client, including
// older ones with a different catalog, can overwrite it.
func FetchCatalog() (Catalog, error) {
	out, err := exec.Command("codex", "debug", "models").Output()
	if err != nil {
		return Catalog{}, fmt.Errorf("codex debug models: %w", err)
	}
	return ParseCatalog(out)
}

// ParseCatalog parses model catalog JSON as printed by codex debug models.
func ParseCatalog(data []byte) (Catalog, error) {
	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return Catalog{}, fmt.Errorf("parse model catalog: %w", err)
	}
	return catalog, nil
}

// Latest returns the newest listed point release in slug's family that
// supports effort. It returns slug itself when nothing newer qualifies.
func (c Catalog) Latest(slug, effort string) string {
	major, current, variant, ok := parse(slug)
	if !ok {
		return slug
	}
	best, bestVersion := slug, current
	for _, m := range c.Models {
		if m.Visibility != "list" || !supports(m, effort) {
			continue
		}
		mMajor, version, mVariant, ok := parse(m.Slug)
		if !ok || mMajor != major || mVariant != variant {
			continue
		}
		if compare(version, bestVersion) > 0 {
			best, bestVersion = m.Slug, version
		}
	}
	return best
}

func parse(slug string) (major string, version []int, variant string, ok bool) {
	match := slugRE.FindStringSubmatch(slug)
	if match == nil {
		return "", nil, "", false
	}
	if match[2] != "" {
		for _, part := range strings.Split(match[2], ".") {
			n, _ := strconv.Atoi(part)
			version = append(version, n)
		}
	}
	return match[1], version, match[3], true
}

// compare orders point versions, treating missing parts as zero so that
// gpt-6-sol (no point version) sorts before gpt-6.1-sol.
func compare(a, b []int) int {
	for i := 0; i < len(a) || i < len(b); i++ {
		var x, y int
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

func supports(m Model, effort string) bool {
	for _, level := range m.ReasoningLevels {
		if level.Effort == effort {
			return true
		}
	}
	return false
}
