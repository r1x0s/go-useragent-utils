package useragent

import (
	"strconv"
	"strings"
)

// Option defines a function to configure the generation process.
type Option func(*generateOptions)

type generateOptions struct {
	browser    BrowserName
	os         OSName
	minVersion Version
	maxVersion Version
	withWeight bool // If true, newer versions are more likely to be picked

	// Header options
	withSecCHUA            bool
	withSecCHUAFullVersion bool
	withSecCHUAPlatform    bool
	withSecCHUAPlatformVer bool
	withSecCHUAMobile      bool
	withSecCHUABitness     bool
	withSecCHUAArch        bool
	withSecCHUAFormFactors bool
	withSecCHUAModel       bool
	withSecCHUAWow64       bool
}

// defaultOptions returns the default configuration.
func defaultOptions() *generateOptions {
	return &generateOptions{
		browser:    Chrome,
		os:         Windows,
		withWeight: true, // Default to weighted selection
	}
}

// WithBrowser sets the target browser.
func WithBrowser(b BrowserName) Option {
	return func(o *generateOptions) {
		o.browser = b
	}
}

// WithOS sets the target operating system.
func WithOS(os OSName) Option {
	return func(o *generateOptions) {
		o.os = os
	}
}

// WithMinVersion sets the minimum allowed version.
// Accepts string like "133.0.0.0" or "145.2".
func WithMinVersion(v string) Option {
	return func(o *generateOptions) {
		o.minVersion = parseVersionString(v)
	}
}

// WithMaxVersion sets the maximum allowed version.
func WithMaxVersion(v string) Option {
	return func(o *generateOptions) {
		o.maxVersion = parseVersionString(v)
	}
}

// WithMinVersionStruct sets the minimum allowed version using a Version struct.
func WithMinVersionStruct(v Version) Option {
	return func(o *generateOptions) {
		o.minVersion = v
	}
}

// WithMaxVersionStruct sets the maximum allowed version using a Version struct.
func WithMaxVersionStruct(v Version) Option {
	return func(o *generateOptions) {
		o.maxVersion = v
	}
}

// WithWeightedSelection enables or disables weighted random selection (favoring newer versions).
func WithWeightedSelection(enable bool) Option {
	return func(o *generateOptions) {
		o.withWeight = enable
	}
}

// WithClientHints enables generation of standard Client Hints headers.
func WithClientHints() Option {
	return func(o *generateOptions) {
		o.withSecCHUA = true
		o.withSecCHUAMobile = true
		o.withSecCHUAPlatform = true
	}
}

// WithAllClientHints enables generation of all available Client Hints headers.
func WithAllClientHints() Option {
	return func(o *generateOptions) {
		o.withSecCHUA = true
		o.withSecCHUAFullVersion = true
		o.withSecCHUAPlatform = true
		o.withSecCHUAPlatformVer = true
		o.withSecCHUAMobile = true
		o.withSecCHUABitness = true
		o.withSecCHUAArch = true
		o.withSecCHUAFormFactors = true
		o.withSecCHUAModel = true
		o.withSecCHUAWow64 = true
	}
}

func parseVersionString(s string) Version {
	parts := strings.Split(s, ".")
	var components []int
	for _, p := range parts {
		if i, err := strconv.Atoi(p); err == nil {
			components = append(components, i)
		} else {
			// Stop parsing on first non-digit
			break
		}
	}
	return Version{Components: components}
}
