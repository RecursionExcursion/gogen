# <img src="https://golang.org/favicon.ico" width="24" alt="Go logo"> gogen (v1.0.0)

[![Go Reference](https://pkg.go.dev/badge/github.com/RecursionExcursion/gogen.svg)](https://pkg.go.dev/github.com/RecursionExcursion/gogen)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**gogen** is a lightweight Go code generator that compiles small Go programs on the fly — ideal for creating quick, single-purpose binaries (e.g., shell scripts, workspace openers, or URL launchers).

---

## Features

- Generate and build temporary Go binaries on the fly  
- Supports Windows, macOS, and Linux targets  
- Includes cleanup function to remove temp build artifacts  
- Logger injection for full caller control  
- CGO disabled by default for portability  

---

## Usage

```go
params := gogen.CreateExeParams{
	Name:     "open-workspace",
	Arch:     gogen.internal.Linux64, // or Win64, Mac64
	Commands: []string{
		`cmd:code ~/projects/foo`,
		`url:https://github.com`,
	},
	Logger: log.New(os.Stdout, "[gogen] ", log.LstdFlags),
}

binPath, exeName, cleanup, err := gogen.GenerateGoExe(params)
if err != nil {
	log.Fatal(err)
}
defer cleanup()

// Serve binary or run it
fmt.Println("Binary ready at:", binPath)
```

## Logging

You can pass a custom *log.Logger via params.Logger.

If you pass nil, logging is silently discarded via io.Discard.

## Cleanup

Call the returned cleanup() after you're done using the generated binary. It removes the temp build directory.

```go
defer cleanup() // Important!
```

## ⚠️ Cross-Compilation Requirements

To build binaries for other platforms, Go must have the standard library installed for each target. Run:

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go install std
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go install std
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go install std

   Without this, cross-compilation may fail with "missing std" errors.

## Installation

go get github.com/RecursionExcursion/gogen

## License

MIT © RecursionExcursion
