package main

import (
	"fmt"

	useragent "github.com/r1x0s/go-useragent-utils/generator"
)

func main() {
	// Initialize generator
	g, _ := useragent.New()
	// Generate with options
	res, _ := g.Generate(
		useragent.WithBrowser(useragent.Chrome),
		useragent.WithOS(useragent.Windows),
		useragent.WithAllClientHints(),
		useragent.WithMinVersion("136"), // Variable length support
	)
	fmt.Println("UA:", res.UserAgent)
	fmt.Println("Headers:", res.Headers)
}
