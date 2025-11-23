package useragent

import (
	_ "embed"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

//go:embed browsers.yaml
var browsersYAML []byte

// browserData holds the flattened data for internal use.
type browserData struct {
	versions   []Version
	uaTemplate string
}

// dataStore holds all loaded browser data.
type dataStore struct {
	data map[BrowserName]map[OSName]*browserData
}

// loadData parses the embedded YAML and returns a structured data store.
func loadData() (*dataStore, error) {
	var config Config
	if err := yaml.Unmarshal(browsersYAML, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedded data: %w", err)
	}

	store := &dataStore{
		data: make(map[BrowserName]map[OSName]*browserData),
	}

	for browserStr, platforms := range config.Browsers {
		browser := BrowserName(browserStr)
		store.data[browser] = make(map[OSName]*browserData)

		for osStr, pConfig := range platforms {
			osName := OSName(osStr)

			bd := &browserData{
				uaTemplate: pConfig.UATemplate,
				versions:   make([]Version, 0),
			}

			// Recursively parse versions
			bd.versions = parseVersions([]int{}, pConfig.Versions)

			// Sort versions descending (newest first)
			sort.Slice(bd.versions, func(i, j int) bool {
				return bd.versions[i].Compare(bd.versions[j]) > 0
			})

			store.data[browser][osName] = bd
		}
	}

	return store, nil
}

func parseVersions(prefix []int, current interface{}) []Version {
	var results []Version

	// Helper to copy prefix
	makePrefix := func(p []int, next int) []int {
		newP := make([]int, len(p)+1)
		copy(newP, p)
		newP[len(p)] = next
		return newP
	}

	switch v := current.(type) {
	case map[int]interface{}:
		for key, val := range v {
			results = append(results, parseVersions(makePrefix(prefix, key), val)...)
		}
	case map[interface{}]interface{}:
		// Handle case where yaml unmarshals keys as interface{}
		for key, val := range v {
			if keyInt, ok := toInt(key); ok {
				results = append(results, parseVersions(makePrefix(prefix, keyInt), val)...)
			}
		}
	case []interface{}:
		// Leaf list of numbers
		for _, item := range v {
			if itemInt, ok := toInt(item); ok {
				results = append(results, Version{Components: makePrefix(prefix, itemInt)})
			}
		}
	case []int:
		// Leaf list of ints (if unmarshaled directly)
		for _, item := range v {
			results = append(results, Version{Components: makePrefix(prefix, item)})
		}
	case nil:
		// End of chain, current prefix is the version
		if len(prefix) > 0 {
			results = append(results, Version{Components: prefix})
		}
	}

	return results
}

func toInt(i interface{}) (int, bool) {
	switch v := i.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}
