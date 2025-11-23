# 🌐 go-useragent-utils

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A powerful and flexible Go library for generating realistic User-Agent strings and Client Hints headers. Perfect for web scraping, testing, and automation tasks.

> **Repository**: [github.com/r1x0s/go-useragent-utils](https://github.com/r1x0s/go-useragent-utils)

## ✨ Features

### Current Implementation

- ✅ **Desktop Chrome Support** (Windows, versions 133+)
- ✅ **Variable Version Length** - Support for any version format (`133`, `133.0`, `133.0.6943.53`)
- ✅ **Flexible Filtering** - Filter by browser, OS, min/max version
- ✅ **Weighted Random Selection** - Newer versions are selected more frequently
- ✅ **Complete Client Hints Support**:
  - `Sec-CH-UA`
  - `Sec-CH-UA-Full-Version-List`
  - `Sec-CH-UA-Platform`
  - `Sec-CH-UA-Platform-Version`
  - `Sec-CH-UA-Mobile`
  - `Sec-CH-UA-Bitness`
  - `Sec-CH-UA-Arch`
  - `Sec-CH-UA-Form-Factors`
  - `Sec-CH-UA-Model`
  - `Sec-CH-UA-Wow64`
- ✅ **GREASE Support** - Automatic randomized GREASE brands for realistic headers
- ✅ **Auto-Update Tool** - Fetch latest Chrome versions from official sources
- ✅ **Zero Dependencies** (runtime) - Embedded YAML data, no external files needed
- ✅ **Type-Safe API** - Functional Options pattern for clean configuration

### 🚀 Planned Features

- 🔜 **Mobile Platform Support** (Android, iOS)
- 🔜 **Additional OS Support** (macOS, Linux)
- 🔜 **Multi-Browser Support**:
  - Firefox
  - Safari
  - Edge
- 🔜 **Custom User-Agent Templates**
- 🔜 **Version History Management**
- 🔜 **Fingerprint Consistency** - Generate matching headers for the same session

## 📦 Installation

```bash
go get github.com/r1x0s/go-useragent-utils
```

## 🎯 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/r1x0s/go-useragent-utils/pkg/useragent"
)

func main() {
    // Initialize generator
    gen, err := useragent.New()
    if err != nil {
        panic(err)
    }

    // Generate User-Agent with default settings
    result, err := gen.Generate()
    if err != nil {
        panic(err)
    }

    fmt.Println("User-Agent:", result.UserAgent)
    // Output: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 
    //         (KHTML, like Gecko) Chrome/133.0.6943.126 Safari/537.36
}
```

### With Client Hints Headers

```go
result, err := gen.Generate(
    useragent.WithClientHints(),
)

for key, value := range result.Headers {
    fmt.Printf("%s: %s\n", key, value)
}
// Output:
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...
// Sec-CH-UA: "Not?A_Brand";v="99", "Google Chrome";v="133", "Chromium";v="133"
// Sec-CH-UA-Mobile: ?0
// Sec-CH-UA-Platform: "Windows"
```

### Advanced Filtering

```go
result, err := gen.Generate(
    useragent.WithBrowser(useragent.Chrome),
    useragent.WithOS(useragent.Windows),
    useragent.WithMinVersion("133.0"),
    useragent.WithMaxVersion("134.0"),
    useragent.WithAllClientHints(),
    useragent.WithWeightedSelection(true),
)
```

### All Available Options

```go
// Browser and OS selection
useragent.WithBrowser(useragent.Chrome)
useragent.WithOS(useragent.Windows)

// Version filtering (supports variable length: "133", "133.0", "133.0.6943.53")
useragent.WithMinVersion("133.0")
useragent.WithMaxVersion("134.0")
useragent.WithMinVersionStruct(useragent.Version{Components: []int{133, 0}})
useragent.WithMaxVersionStruct(useragent.Version{Components: []int{134, 0}})

// Selection strategy
useragent.WithWeightedSelection(true)  // Favor newer versions (default: true)

// Client Hints headers
useragent.WithClientHints()     // Standard headers (UA, Mobile, Platform)
useragent.WithAllClientHints()  // All available headers
```

## 🔧 Updating Chrome Versions

The library includes a tool to automatically fetch and update Chrome versions:

```bash
go run ./cmd/update-data
```

This will:
1. Fetch the latest Chrome versions from Google Chrome Labs API
2. Filter versions >= 133
3. Update `data/browsers.yaml` with new versions
4. Preserve existing data structure

After updating, rebuild your application to embed the new data.

## 📁 Project Structure

```
go-useragent-utils/
├── cmd/
│   └── example/          # Example usage
│       └── main.go
│   └── update-data/          # Auto-update tool
│       └── main.go
├── data/
│   └── browsers.yaml         # Version database (embedded)
├── generator/
│   ├── types.go          # Core types and constants
│   ├── data.go           # YAML loading and parsing
│   ├── generator.go      # Main generation logic
│   ├── headers.go        # Client Hints generation
│   ├── options.go        # Functional options
│   └── browsers.yaml     # Embedded copy of data
├── README.md
├── go.mod
└── go.sum
```

## 🧪 Testing

```bash
# Run all tests
go test -v ./generator

# Run with coverage
go test -cover ./generator
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Chrome version data from [Google Chrome Labs](https://googlechromelabs.github.io/chrome-for-testing/)
- Inspired by the need for realistic browser fingerprinting in automation

---

**Made with ❤️ by [r1x0s](https://github.com/r1x0s)**
