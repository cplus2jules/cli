package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spicetify/cli/resources/selectors"
)

type SelectorMap struct {
	Version             string            `json:"version"`
	SpotifyVersionRange string            `json:"spotifyVersionRange"`
	Selectors           map[string]string `json:"selectors"`
	CSSVarAliases       map[string]string `json:"cssVarAliases"`
}

func LoadSelectorMap(spotifyVersion string) (*SelectorMap, error) {
	if spotifyVersion == "" {
		PrintWarning("Empty Spotify version provided, using fallback selector map")
		return loadFallbackMap()
	}

	parts := strings.Split(spotifyVersion, ".")
	if len(parts) < 3 {
		PrintWarning("Invalid Spotify version format (" + spotifyVersion + "), using fallback selector map")
		return loadFallbackMap()
	}
	shortVersion := strings.Join(parts[:3], ".")

	// Try exact match first
	filename := fmt.Sprintf("%s.json", shortVersion)
	data, err := selectorassets.FS.ReadFile(filename)
	if err == nil {
		return parseSelectorMap(data)
	}

	// Try to find the closest version dynamically if exact doesn't exist
	entries, _ := selectorassets.FS.ReadDir(".")
	bestMatch := ""
	for _, entry := range entries {
		if entry.Name() == "fallback.json" || entry.Name() == "selectorassets.go" {
			continue
		}
		// A simple strategy is falling back to minor matches, or any other closest match.
		// For now we just check if it shares the first 2 segments (e.g. 1.2)
		entryParts := strings.Split(entry.Name(), ".")
		if len(entryParts) >= 3 && parts[0] == entryParts[0] && parts[1] == entryParts[1] {
			// Find highest patch version that is <= user's patch version?
			// To keep it robust, we select the first one we find for that minor block.
			bestMatch = entry.Name()
			break // in reality we might want a stricter semver parsing loop
		}
	}

	if bestMatch != "" {
		data, err = selectorassets.FS.ReadFile(bestMatch)
		if err == nil {
			PrintInfo("Using selector map from " + bestMatch + " (closest match for " + shortVersion + ")")
			return parseSelectorMap(data)
		}
	}

	PrintWarning("No selector map found for Spotify " + spotifyVersion + ", using fallback")
	return loadFallbackMap()
}

func loadFallbackMap() (*SelectorMap, error) {
	data, err := selectorassets.FS.ReadFile("fallback.json")
	if err != nil {
		return nil, fmt.Errorf("critical: fallback selector map missing from binary")
	}
	return parseSelectorMap(data)
}

func parseSelectorMap(data []byte) (*SelectorMap, error) {
	var sm SelectorMap
	if err := json.Unmarshal(data, &sm); err != nil {
		return nil, fmt.Errorf("failed to parse selector map: %w", err)
	}
	return &sm, nil
}
