// Package internal holds package-level constants and shared types for tamcp.
package internal

import "runtime"

// NAME is the program identity (binary name, service short name, log tag).
const NAME = "tamcp"

// Version, Commit, BuildDate are wired in via -ldflags from the Makefile.
var (
	Version   = "0.0.1"
	Commit    = "dev"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)
