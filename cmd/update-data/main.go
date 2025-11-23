package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	chromeVersionsURL = "https://googlechromelabs.github.io/chrome-for-testing/known-good-versions.json"
	minMajorVersion   = 133
	dataFile          = "data/browsers.yaml"
)

type ChromeVersionsResponse struct {
	Versions []struct {
		Version string `json:"version"`
	} `json:"versions"`
}

// Config matches the structure in pkg/useragent/types.go but simplified for manipulation
type Config struct {
	Browsers map[string]map[string]PlatformConfig `yaml:"browsers"`
}

type PlatformConfig struct {
	UATemplate string              `yaml:"ua_template"`
	Versions   map[int]interface{} `yaml:"versions"`
}

func main() {
	fmt.Println("Fetching Chrome versions...")
	versions, err := fetchChromeVersions()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d versions. Filtering for >= %d...\n", len(versions), minMajorVersion)

	filtered := filterVersions(versions)
	fmt.Printf("Kept %d versions.\n", len(filtered))

	fmt.Println("Updating data/browsers.yaml...")
	if err := updateYAML(filtered); err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}

func fetchChromeVersions() ([]string, error) {
	resp, err := http.Get(chromeVersionsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data ChromeVersionsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var versions []string
	for _, v := range data.Versions {
		versions = append(versions, v.Version)
	}
	return versions, nil
}

func filterVersions(versions []string) []string {
	var res []string
	for _, v := range versions {
		parts := strings.Split(v, ".")
		if len(parts) > 0 {
			major, err := strconv.Atoi(parts[0])
			if err == nil && major >= minMajorVersion {
				res = append(res, v)
			}
		}
	}
	return res
}

func updateYAML(newVersions []string) error {
	// Read existing file
	path, _ := filepath.Abs(dataFile)
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	if err := yaml.Unmarshal(content, &config); err != nil {
		return err
	}

	// Initialize if empty
	if config.Browsers == nil {
		config.Browsers = make(map[string]map[string]PlatformConfig)
	}
	if config.Browsers["chrome"] == nil {
		config.Browsers["chrome"] = make(map[string]PlatformConfig)
	}

	// Update Windows versions
	winConfig, ok := config.Browsers["chrome"]["windows"]
	if !ok {
		// Should exist based on our seed, but handle anyway
		winConfig = PlatformConfig{
			UATemplate: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/{{version}} Safari/537.36",
			Versions:   make(map[int]interface{}),
		}
	}

	// Merge new versions into the map
	for _, vStr := range newVersions {
		parts := parseVersion(vStr)
		if len(parts) == 0 {
			continue
		}
		addToMap(winConfig.Versions, parts)
	}

	config.Browsers["chrome"]["windows"] = winConfig

	// Write back
	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}

func parseVersion(v string) []int {
	parts := strings.Split(v, ".")
	var res []int
	for _, p := range parts {
		i, _ := strconv.Atoi(p)
		res = append(res, i)
	}
	return res
}

// addToMap recursively adds the version components to the nested map
func addToMap(m map[int]interface{}, components []int) {
	if len(components) == 0 {
		return
	}
	head := components[0]
	tail := components[1:]

	if len(tail) == 0 {
		// Leaf node, but in our structure leaves are usually in a list at the end.
		// However, to keep it simple and consistent with the recursive structure:
		// If we are at the end, we ensure the key exists.
		// But wait, our structure is:
		// 133:
		//   0:
		//     6943: [53, 98]
		// So the last component is in a list.
		// Let's adjust logic:
		// If len(tail) == 0, we are at the leaf. But we can't add a key to a map and have no value?
		// Actually, the previous level should have handled this.
		// Let's change strategy:
		// We traverse until the LAST component.
		return
	}

	// If we have 1 element left in tail, that element is the leaf value to be added to the list.
	if len(tail) == 1 {
		// Check if m[head] exists
		if _, ok := m[head]; !ok {
			m[head] = []int{}
		}

		// It could be a map (if we have deeper versions) or a list (if we are at leaf).
		// Our schema assumes fixed depth? No, variable.
		// If it's a list, we append.
		// If it's a map, we have a conflict? (e.g. 133.0 vs 133.0.1)
		// For Chrome it's always 4 parts.
		// Let's handle the list case.

		switch v := m[head].(type) {
		case []int:
			// Check if exists
			exists := false
			for _, x := range v {
				if x == tail[0] {
					exists = true
					break
				}
			}
			if !exists {
				// Append and sort
				v = append(v, tail[0])
				sort.Ints(v)
				m[head] = v
			}
		case []interface{}:
			// YAML unmarshal might give []interface{}
			var ints []int
			for _, x := range v {
				if i, ok := x.(int); ok {
					ints = append(ints, i)
				}
			}
			exists := false
			for _, x := range ints {
				if x == tail[0] {
					exists = true
					break
				}
			}
			if !exists {
				ints = append(ints, tail[0])
				sort.Ints(ints)
				m[head] = ints
			}
		case map[int]interface{}:
			// We are trying to add a leaf but there is a map?
			// This means we have 133.0.1 (map) and we want to add 133.0 (leaf)?
			// Chrome versions are consistent depth usually.
			// Let's assume for now we just recurse if it's a map.
			addToMap(v, tail)
		case map[interface{}]interface{}:
			// Convert to map[int]interface{}
			newMap := make(map[int]interface{})
			for k, val := range v {
				if kInt, ok := k.(int); ok {
					newMap[kInt] = val
				}
			}
			m[head] = newMap
			addToMap(newMap, tail)
		case nil:
			// Create list
			m[head] = []int{tail[0]}
		}
		return
	}

	// More than 1 element in tail, so we need a map at this level
	if _, ok := m[head]; !ok {
		m[head] = make(map[int]interface{})
	}

	// Ensure it is a map
	switch v := m[head].(type) {
	case map[int]interface{}:
		addToMap(v, tail)
	case map[interface{}]interface{}:
		newMap := make(map[int]interface{})
		for k, val := range v {
			if kInt, ok := k.(int); ok {
				newMap[kInt] = val
			}
		}
		m[head] = newMap
		addToMap(newMap, tail)
	case []int:
		// We have a list (leafs) but need to go deeper?
		// This implies mixed depth. Convert list to map?
		// E.g. we had 133.0 -> [1, 2]
		// Now we want 133.0.1.5
		// This is tricky. For Chrome it shouldn't happen.
		// Let's ignore for now or overwrite?
		// Let's just return to avoid panic.
		return
	case nil:
		newMap := make(map[int]interface{})
		m[head] = newMap
		addToMap(newMap, tail)
	}
}
