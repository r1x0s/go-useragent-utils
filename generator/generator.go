package useragent

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Generator is the main entry point for generating user agents.
type Generator struct {
	store *dataStore
	rng   *rand.Rand
}

// New creates a new Generator with loaded data.
func New() (*Generator, error) {
	store, err := loadData()
	if err != nil {
		return nil, err
	}
	return &Generator{
		store: store,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Result contains the generated User-Agent string and headers.
type Result struct {
	UserAgent string
	Headers   map[string]string
}

// Generate creates a new User-Agent and optional headers based on the provided options.
func (g *Generator) Generate(opts ...Option) (*Result, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 1. Get browser data
	platforms, ok := g.store.data[options.browser]
	if !ok {
		return nil, fmt.Errorf("browser %s not found", options.browser)
	}
	bd, ok := platforms[options.os]
	if !ok {
		return nil, fmt.Errorf("os %s not found for browser %s", options.os, options.browser)
	}

	// 2. Filter versions
	candidates := g.filterVersions(bd.versions, options)
	if len(candidates) == 0 {
		return nil, errors.New("no versions found matching criteria")
	}

	// 3. Select version
	selectedVer := g.selectVersion(candidates, options.withWeight)

	// 4. Build User-Agent string
	ua := strings.ReplaceAll(bd.uaTemplate, "{{version}}", selectedVer.String())

	// 5. Build Headers
	headers := g.generateHeaders(selectedVer, options)
	headers["User-Agent"] = ua

	return &Result{
		UserAgent: ua,
		Headers:   headers,
	}, nil
}

func (g *Generator) filterVersions(versions []Version, opts *generateOptions) []Version {
	var filtered []Version
	for _, v := range versions {
		// Min version check
		if len(opts.minVersion.Components) > 0 && v.Compare(opts.minVersion) < 0 {
			continue
		}
		// Max version check
		if len(opts.maxVersion.Components) > 0 && v.Compare(opts.maxVersion) > 0 {
			continue
		}
		filtered = append(filtered, v)
	}
	return filtered
}

func (g *Generator) selectVersion(versions []Version, weighted bool) Version {
	if len(versions) == 0 {
		return Version{} // Should not happen if filtered correctly
	}
	if !weighted || len(versions) == 1 {
		return versions[g.rng.Intn(len(versions))]
	}

	// Weighted random selection: newer versions have higher weight.
	// Simple linear weight: index 0 (newest) gets weight N, index N-1 gets weight 1.
	// Or maybe exponential? Let's stick to linear for now as requested "more new versions".
	// Since versions are sorted descending (newest at 0), we give higher weight to lower indices.

	n := len(versions)
	totalWeight := (n * (n + 1)) / 2 // Sum of 1..n
	r := g.rng.Intn(totalWeight)

	// Find which item corresponds to r
	// Weights: [n, n-1, ..., 1]
	currentWeight := 0
	for i, v := range versions {
		weight := n - i
		currentWeight += weight
		if r < currentWeight {
			return v
		}
	}

	return versions[0] // Fallback
}
