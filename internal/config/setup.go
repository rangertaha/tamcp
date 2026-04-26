package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rangertaha/tamcp/internal"
)

var (
	//go:embed config.hcl
	defaultConfig []byte

	DefaultConfigFileName         = "config.hcl"
	DefaultGlobalConfigDir        = "/etc"
	DefaultGlobalDataDir          = "/var/lib"
	DefaultUserConfigDir          = ".config"
	DefaultUserDataDir            = ".local/share"
	DefaultSystemdServiceDir      = "/etc/systemd/system"
	DefaultSystemdServiceFileName = "tamcp.service"
	DefaultSystemdServiceTemplate = `[Unit]
Description=tamcp MCP server (technical-analysis indicators)
After=network.target

[Service]
ExecStart=/usr/bin/tamcp server

[Install]
WantedBy=multi-user.target
`
	GlobalConfigDir string
	GlobalDataDir   string
	UserConfigDir   string
	UserDataDir     string
)

func init() {
	programName := strings.ToLower(internal.NAME)
	if cfgDir, err := os.UserConfigDir(); err == nil {
		DefaultUserConfigDir = cfgDir
	} else if home, hErr := os.UserHomeDir(); hErr == nil {
		DefaultUserConfigDir = filepath.Join(home, ".config")
	} else {
		DefaultUserConfigDir = "~/.config"
	}

	GlobalConfigDir = filepath.Join(DefaultGlobalConfigDir, programName)
	GlobalDataDir = filepath.Join(DefaultGlobalDataDir, programName)
	UserConfigDir = filepath.Join(DefaultUserConfigDir, programName)
	UserDataDir = filepath.Join(DefaultUserDataDir, programName)
}

// CreateConfigFile writes the default config to dir/filename, falling back
// to sudo when the destination requires elevated privileges.
func CreateConfigFile(dir, filename string) error {
	if strings.TrimSpace(dir) == "" {
		return errors.New("config dir is required")
	}
	if strings.TrimSpace(filename) == "" {
		return errors.New("config filename is required")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		mkdirCmd := exec.Command("sudo", "mkdir", "-p", dir)
		mkdirCmd.Stdin = os.Stdin
		mkdirCmd.Stdout = os.Stdout
		mkdirCmd.Stderr = os.Stderr
		if sudoErr := mkdirCmd.Run(); sudoErr != nil {
			return sudoErr
		}
	}

	file := filepath.Join(dir, filename)
	if _, err := os.Stat(file); err == nil {
		return errors.New("file already exists: " + file)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.WriteFile(file, defaultConfig, 0644); err == nil {
		fmt.Println("config file created: " + file)
		return nil
	} else if !os.IsPermission(err) {
		return err
	}

	teeCmd := exec.Command("sudo", "tee", file)
	teeCmd.Stdin = bytes.NewReader(defaultConfig)
	teeCmd.Stdout = io.Discard
	teeCmd.Stderr = os.Stderr
	if err := teeCmd.Run(); err != nil {
		return err
	}
	fmt.Println("config file created: " + file)
	return nil
}

func CreateConfigPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("config path is required")
	}
	if _, err := os.Stat(path); err == nil {
		return errors.New("path already exists: " + path)
	}
	if !strings.HasSuffix(path, ".hcl") {
		return errors.New("file name must end with .hcl")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return CreateConfigFile(dir, filepath.Base(path))
}

func RemoveConfigPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("config path is required")
	}
	if !strings.HasSuffix(path, ".hcl") {
		return errors.New("file name must end with .hcl")
	}
	return RemoveConfigFile(filepath.Dir(path), filepath.Base(path))
}

// RemoveConfigFile removes a config file and its parent directory if empty.
func RemoveConfigFile(dir, filename string) error {
	file := filepath.Join(dir, filename)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(file); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		rmCmd := exec.Command("sudo", "rm", file)
		rmCmd.Stdin = os.Stdin
		rmCmd.Stdout = os.Stdout
		rmCmd.Stderr = os.Stderr
		if sudoErr := rmCmd.Run(); sudoErr != nil {
			return sudoErr
		}
	}
	fmt.Println("Removed: " + file)

	entries, err := os.ReadDir(dir)
	if err == nil && len(entries) == 0 {
		if rmErr := os.Remove(dir); rmErr != nil && os.IsPermission(rmErr) {
			rmDirCmd := exec.Command("sudo", "rmdir", dir)
			rmDirCmd.Stdin = os.Stdin
			rmDirCmd.Stdout = os.Stdout
			rmDirCmd.Stderr = os.Stderr
			_ = rmDirCmd.Run()
		}
	}
	return nil
}

// CreateSystemdServiceFile installs and enables the systemd unit, falling back to sudo.
func CreateSystemdServiceFile() error {
	if err := os.MkdirAll(DefaultSystemdServiceDir, 0755); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		mkdirCmd := exec.Command("sudo", "mkdir", "-p", DefaultSystemdServiceDir)
		mkdirCmd.Stdin = os.Stdin
		mkdirCmd.Stdout = os.Stdout
		mkdirCmd.Stderr = os.Stderr
		if sudoErr := mkdirCmd.Run(); sudoErr != nil {
			return sudoErr
		}
	}

	file := filepath.Join(DefaultSystemdServiceDir, DefaultSystemdServiceFileName)
	if _, err := os.Stat(file); err == nil {
		return errors.New("file already exists: " + file)
	} else if !os.IsNotExist(err) {
		return err
	}

	content := []byte(DefaultSystemdServiceTemplate)
	if err := os.WriteFile(file, content, 0644); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		teeCmd := exec.Command("sudo", "tee", file)
		teeCmd.Stdin = bytes.NewReader(content)
		teeCmd.Stdout = io.Discard
		teeCmd.Stderr = os.Stderr
		if sudoErr := teeCmd.Run(); sudoErr != nil {
			return sudoErr
		}
	}

	for _, args := range [][]string{
		{"sudo", "systemctl", "enable", DefaultSystemdServiceFileName},
		{"sudo", "systemctl", "start", DefaultSystemdServiceFileName},
	} {
		c := exec.Command(args[0], args[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return err
		}
	}
	fmt.Println("systemd service file created: " + file)
	return nil
}

// RemoveSystemdServiceFile stops, disables, and removes the systemd unit.
func RemoveSystemdServiceFile() error {
	file := filepath.Join(DefaultSystemdServiceDir, DefaultSystemdServiceFileName)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}
	for _, args := range [][]string{
		{"sudo", "systemctl", "stop", DefaultSystemdServiceFileName},
		{"sudo", "systemctl", "disable", DefaultSystemdServiceFileName},
	} {
		c := exec.Command(args[0], args[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		_ = c.Run()
	}
	if err := os.Remove(file); err != nil {
		if !os.IsPermission(err) {
			return err
		}
		rmCmd := exec.Command("sudo", "rm", file)
		rmCmd.Stdin = os.Stdin
		rmCmd.Stdout = os.Stdout
		rmCmd.Stderr = os.Stderr
		if sudoErr := rmCmd.Run(); sudoErr != nil {
			return sudoErr
		}
	}
	fmt.Println("Removed: " + file)
	return nil
}
