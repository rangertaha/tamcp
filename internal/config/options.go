package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Option configures Config construction.
type Option func(*Config) error

func WithDebug(debug bool) Option {
	return func(c *Config) error {
		c.Debug = debug
		return nil
	}
}

func WithDataDir(dataDir string) Option {
	return func(c *Config) error {
		if strings.TrimSpace(dataDir) != "" {
			c.DataDir = dataDir
		}
		return nil
	}
}

func WithLogFile(file string) Option {
	return func(c *Config) error {
		if strings.TrimSpace(file) != "" {
			c.Logging.File = file
		}
		return nil
	}
}

func WithLogLevel(level string) Option {
	return func(c *Config) error {
		if strings.TrimSpace(level) != "" {
			c.Logging.Level = level
		}
		return nil
	}
}

func WithConfig(cfgs ...*Config) Option {
	return func(c *Config) error {
		for _, cfg := range cfgs {
			if cfg == nil {
				continue
			}
			c.Merge(cfg)
		}
		return nil
	}
}

func WithConfigFile(file string) Option {
	return func(c *Config) error {
		if strings.TrimSpace(file) == "" {
			return nil
		}
		cfg, err := Load(file)
		if err != nil {
			return err
		}
		c.Merge(cfg)
		return nil
	}
}

// WithGlobalConfigFile merges /etc/tamcp/config.hcl if it exists. Missing is not an error.
func WithGlobalConfigFile() Option {
	return withOptionalConfigFile(filepath.Join(GlobalConfigDir, DefaultConfigFileName))
}

// WithUserConfigFile merges ~/.config/tamcp/config.hcl if it exists. Missing is not an error.
func WithUserConfigFile() Option {
	return withOptionalConfigFile(filepath.Join(UserConfigDir, DefaultConfigFileName))
}

func withOptionalConfigFile(path string) Option {
	return func(c *Config) error {
		if strings.TrimSpace(path) == "" {
			return nil
		}
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		cfg, err := Load(path)
		if err != nil {
			return err
		}
		c.Merge(cfg)
		return nil
	}
}
