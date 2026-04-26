//go:build !windows

// Package winservice provides Windows-Service stubs on non-Windows builds so
// the rest of the codebase compiles. All operations return errUnsupported.
package winservice

import (
	"context"
	"errors"
)

const (
	Name        = "tamcp"
	DisplayName = "tamcp"
	Description = "MCP server for technical-analysis indicators"
)

var errUnsupported = errors.New("windows service operations are only supported on windows")

type Handler func(ctx context.Context) error

func IsService() (bool, error)          { return false, nil }
func Run(run Handler) error             { return errUnsupported }
func Install(args ...string) error      { return errUnsupported }
func Uninstall() error                  { return errUnsupported }
func StartService(args ...string) error { return errUnsupported }
func StopService() error                { return errUnsupported }
