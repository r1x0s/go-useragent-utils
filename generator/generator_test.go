package useragent

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	g, err := New()
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	t.Run("Default", func(t *testing.T) {
		res, err := g.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		if res.UserAgent == "" {
			t.Error("Empty User-Agent")
		}
		if res.Headers["User-Agent"] != res.UserAgent {
			t.Error("User-Agent header missing or mismatch")
		}
		t.Logf("UA: %s", res.UserAgent)
	})

	t.Run("WithFilters", func(t *testing.T) {
		// Test with a partial version string
		res, err := g.Generate(
			WithBrowser(Chrome),
			WithOS(Windows),
			WithMinVersion("133.0"),
		)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		if res.UserAgent == "" {
			t.Error("Empty User-Agent")
		}
	})

	t.Run("WithHeaders", func(t *testing.T) {
		res, err := g.Generate(WithClientHints())
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		if _, ok := res.Headers["Sec-CH-UA"]; !ok {
			t.Error("Missing Sec-CH-UA header")
		}
		if _, ok := res.Headers["Sec-CH-UA-Mobile"]; !ok {
			t.Error("Missing Sec-CH-UA-Mobile header")
		}
		if _, ok := res.Headers["User-Agent"]; !ok {
			t.Error("Missing User-Agent header in Headers map")
		}
		t.Logf("Headers: %v", res.Headers)
	})

	t.Run("VariableLengthVersion", func(t *testing.T) {
		// We need to mock data or ensure we have a variable length version in data.
		// Since we only have Chrome 4-part versions in embedded data,
		// we can't easily test generating a 2-part version without modifying the embedded file
		// or mocking the store.
		// However, we can test the Version struct logic directly.
		v1 := Version{Components: []int{145, 2}}
		v2 := Version{Components: []int{145, 2, 0, 0}}

		if v1.String() != "145.2" {
			t.Errorf("Expected 145.2, got %s", v1.String())
		}

		// Compare logic
		if v1.Compare(v2) != 0 {
			// Depending on logic, 145.2 might be equal to 145.2.0.0 or not.
			// My implementation compares up to max length, treating missing as 0.
			// So they should be equal.
			t.Errorf("Expected 145.2 == 145.2.0.0")
		}
	})
	t.Run("VariableLengthVersionInputs", func(t *testing.T) {
		// Test "133"
		res, err := g.Generate(WithMinVersion("133"))
		if err != nil {
			t.Errorf("Failed with '133': %v", err)
		}
		if res == nil || res.UserAgent == "" {
			t.Error("Empty result for '133'")
		}

		// Test "133.0.6943.10000" (should likely fail if we don't have it, or pass if we do)
		// In our data we have 133.0.6943.[53, 98, 126, 141].
		// So 133.0.6943.10000 should result in error (no versions found)
		_, err = g.Generate(WithMinVersion("133.0.6943.10000"))
		if err == nil {
			t.Error("Expected error for very high version, got nil")
		}

		// Test "133.0.0.0"
		res, err = g.Generate(WithMinVersion("133.0.0.0"))
		if err != nil {
			t.Errorf("Failed with '133.0.0.0': %v", err)
		}
	})
}
