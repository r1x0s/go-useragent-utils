package useragent

import (
	"fmt"
	"math/rand"
)

// generateHeaders creates the map of HTTP headers based on options and selected version.
func (g *Generator) generateHeaders(v Version, opts *generateOptions) map[string]string {
	headers := make(map[string]string)

	if opts.withSecCHUA {
		headers["Sec-CH-UA"] = g.formatSecCHUA(v)
	}
	if opts.withSecCHUAMobile {
		headers["Sec-CH-UA-Mobile"] = "?0" // Desktop only for now
	}
	if opts.withSecCHUAPlatform {
		headers["Sec-CH-UA-Platform"] = "\"Windows\""
	}
	if opts.withSecCHUAFullVersion {
		headers["Sec-CH-UA-Full-Version-List"] = g.formatSecCHUAFullVersion(v)
	}
	if opts.withSecCHUAPlatformVer {
		headers["Sec-CH-UA-Platform-Version"] = "\"10.0.0\""
	}
	if opts.withSecCHUABitness {
		headers["Sec-CH-UA-Bitness"] = "\"64\""
	}
	if opts.withSecCHUAArch {
		headers["Sec-CH-UA-Arch"] = "\"x86\""
	}
	if opts.withSecCHUAModel {
		headers["Sec-CH-UA-Model"] = "\"\""
	}
	if opts.withSecCHUAWow64 {
		headers["Sec-CH-UA-Wow64"] = "?0"
	}
	if opts.withSecCHUAFormFactors {
		headers["Sec-CH-UA-Form-Factors"] = "\"Desktop\""
	}

	return headers
}

func (g *Generator) formatSecCHUA(v Version) string {
	grease := g.getGreaseBrand()
	major := 0
	if len(v.Components) > 0 {
		major = v.Components[0]
	}
	return fmt.Sprintf(`"%s";v="99", "Not(A:Brand";v="99", "Google Chrome";v="%d", "Chromium";v="%d"`, grease, major, major)
}

func (g *Generator) formatSecCHUAFullVersion(v Version) string {
	grease := g.getGreaseBrand()
	fullVer := v.String()

	// Ensure full version has at least 4 components for Chrome if possible,
	// but if the version is "145.2", we just output "145.2".
	// However, spec usually expects full version.
	// If it's Chrome, it should be 4 parts.
	// We'll trust the data source or pad if needed?
	// For now, just use v.String() which joins components.

	return fmt.Sprintf(`"%s";v="99.0.0.0", "Not(A:Brand";v="99.0.0.0", "Google Chrome";v="%s", "Chromium";v="%s"`, grease, fullVer, fullVer)
}

func (g *Generator) getGreaseBrand() string {
	brands := []string{"Not(A:Brand", "Not?A_Brand", "Not A;Brand"}
	return brands[rand.Intn(len(brands))]
}
