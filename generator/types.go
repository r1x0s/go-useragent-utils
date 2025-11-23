package useragent

import (
	"fmt"
	"strings"
)

// BrowserName represents the browser name (e.g., "chrome", "firefox").
type BrowserName string

// OSName represents the operating system name (e.g., "windows", "linux").
type OSName string

const (
	Chrome  BrowserName = "chrome"
	Firefox BrowserName = "firefox"
	Safari  BrowserName = "safari"
	Edge    BrowserName = "edge"

	Windows OSName = "windows"
	Linux   OSName = "linux"
	MacOS   OSName = "macos"
	Android OSName = "android"
	IOS     OSName = "ios"
)

// Version represents a semantic version with variable number of components.
type Version struct {
	Components []int
}

func (v Version) String() string {
	parts := make([]string, len(v.Components))
	for i, c := range v.Components {
		parts[i] = fmt.Sprintf("%d", c)
	}
	return strings.Join(parts, ".")
}

// Compare returns -1 if v < other, 1 if v > other, 0 if equal.
func (v Version) Compare(other Version) int {
	len1 := len(v.Components)
	len2 := len(other.Components)
	maxLen := len1
	if len2 > maxLen {
		maxLen = len2
	}

	for i := 0; i < maxLen; i++ {
		v1 := 0
		if i < len1 {
			v1 = v.Components[i]
		}
		v2 := 0
		if i < len2 {
			v2 = other.Components[i]
		}

		if v1 < v2 {
			return -1
		}
		if v1 > v2 {
			return 1
		}
	}
	return 0
}

// Config represents the top-level structure of the YAML file.
type Config struct {
	Browsers map[string]map[string]PlatformConfig `yaml:"browsers"`
}

// PlatformConfig holds the template and version data for a specific OS.
type PlatformConfig struct {
	UATemplate string `yaml:"ua_template"`
	// Versions is a nested map structure.
	// We use map[int]interface{} to support variable depth.
	// The value can be:
	// - map[int]interface{} (next level)
	// - []int (leaf list of patches)
	// - nil (end of version)
	Versions map[int]interface{} `yaml:"versions"`
}
