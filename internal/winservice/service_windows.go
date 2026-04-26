//go:build windows

// Package winservice integrates tamcp with the Windows Service Control Manager.
package winservice

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	Name        = "tamcp"
	DisplayName = "tamcp"
	Description = "MCP server for technical-analysis indicators"
)

// Handler is the service body. It should return when ctx is canceled or the
// underlying transport ends.
type Handler func(ctx context.Context) error

// IsService reports whether the current process was launched by the SCM.
func IsService() (bool, error) { return svc.IsWindowsService() }

type winHandler struct{ run Handler }

func (h *winHandler) Execute(_ []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	const accepts = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- h.run(ctx) }()

	changes <- svc.Status{State: svc.Running, Accepts: accepts}

loop:
	for {
		select {
		case req := <-r:
			switch req.Cmd {
			case svc.Interrogate:
				changes <- req.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				cancel()
				break loop
			}
		case <-done:
			break loop
		}
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
	}
	changes <- svc.Status{State: svc.Stopped}
	return false, 0
}

// Run registers with SCM and blocks until the service exits.
func Run(run Handler) error { return svc.Run(Name, &winHandler{run: run}) }

// Install registers the service with the SCM.
func Install(args ...string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	if existing, err := m.OpenService(Name); err == nil {
		existing.Close()
		return fmt.Errorf("service %q already exists", Name)
	}

	s, err := m.CreateService(Name, exe, mgr.Config{
		DisplayName: DisplayName,
		Description: Description,
		StartType:   mgr.StartAutomatic,
	}, args...)
	if err != nil {
		return err
	}
	s.Close()
	return nil
}

func Uninstall() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(Name)
	if err != nil {
		return fmt.Errorf("service %q not installed: %w", Name, err)
	}
	defer s.Close()

	status, err := s.Query()
	if err == nil && status.State != svc.Stopped {
		_, _ = s.Control(svc.Stop)
		waitForState(s, svc.Stopped, 30*time.Second)
	}
	return s.Delete()
}

func StartService(args ...string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(Name)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.Start(args...)
}

func StopService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(Name)
	if err != nil {
		return err
	}
	defer s.Close()
	if _, err := s.Control(svc.Stop); err != nil {
		return err
	}
	return waitForState(s, svc.Stopped, 30*time.Second)
}

func waitForState(s *mgr.Service, target svc.State, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		status, err := s.Query()
		if err != nil {
			return err
		}
		if status.State == target {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for service state %d", target)
		}
		time.Sleep(300 * time.Millisecond)
	}
}
