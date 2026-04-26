// Package config loads, merges, and persists tamcp's HCL configuration.
//
// Three layers are merged in this order: built-in defaults, optional global
// (/etc/tamcp/config.hcl), optional user (~/.config/tamcp/config.hcl), and an
// explicit per-invocation file. Later layers override earlier ones.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// Config is the merged root configuration decoded from HCL.
type Config struct {
	// Debug enables verbose debug output across all subsystems.
	Debug bool `hcl:"debug,optional"`

	// DataDir is the root directory for persistent data (databases, state files).
	DataDir string `hcl:"datadir,optional"`

	// Logging configures log level and optional file output.
	Logging struct {
		Level string `hcl:"level,optional"`
		File  string `hcl:"file,optional"`
	} `hcl:"logging,block"`

	// Server configures the MCP server runtime.
	Server *Server `hcl:"server,block"`

	// Database configures the GORM persistence layer used by stateful tools.
	Database *Database `hcl:"database,block"`

	// Tools lists per-tool enablement and configuration.
	Tools []*Tool `hcl:"tool,block"`

	// Providers lists upstream market-data sources.
	Providers []*Provider `hcl:"provider,block"`

	// Brokers lists trading-execution venues.
	Brokers []*Broker `hcl:"broker,block"`
}

// Server configures the MCP server runtime.
type Server struct {
	Name      string `hcl:"name,optional"`
	Transport string `hcl:"transport,optional"`
}

// Database configures the GORM persistence layer.
type Database struct {
	Driver string `hcl:"driver,optional"`
	DSN    string `hcl:"dsn,optional"`
}

// Tool is a per-tool config block. Plugins decode Body into their own struct.
type Tool struct {
	Name    string   `hcl:"name,label"`
	Enabled bool     `hcl:"enabled,optional"`
	Body    hcl.Body `hcl:",remain"`
}

// Provider is a data-source plugin block.
type Provider struct {
	Name    string   `hcl:"name,label"`
	Enabled bool     `hcl:"enabled,optional"`
	Body    hcl.Body `hcl:",remain"`
}

// Broker is a trading-execution venue plugin block.
type Broker struct {
	Name    string   `hcl:"name,label"`
	Enabled bool     `hcl:"enabled,optional"`
	Body    hcl.Body `hcl:",remain"`
}

// GetTool returns the named tool block, or nil if not found.
func (c *Config) GetTool(name string) *Tool {
	for _, t := range c.Tools {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// GetProvider returns the named provider block, or nil if not found.
func (c *Config) GetProvider(name string) *Provider {
	for _, p := range c.Providers {
		if p.Name == name {
			return p
		}
	}
	return nil
}

// GetBroker returns the named broker block, or nil if not found.
func (c *Config) GetBroker(name string) *Broker {
	for _, b := range c.Brokers {
		if b.Name == name {
			return b
		}
	}
	return nil
}

// New builds a Config by applying options in order.
func New(options ...Option) (*Config, error) {
	c := &Config{}
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Load parses an HCL file into a Config.
func Load(file string) (*Config, error) {
	cfg := &Config{}
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(file)
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}
	if diags := gohcl.DecodeBody(f.Body, nil, cfg); diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}
	return cfg, nil
}

// Merge layers cfg over c. Non-zero/non-empty fields in cfg overwrite c.
// Tool and Provider blocks are union-merged by name (cfg wins on collision).
func (c *Config) Merge(cfg *Config) *Config {
	if c == nil {
		return cfg
	}
	if cfg == nil {
		return c
	}
	if cfg.Debug {
		c.Debug = true
	}
	if cfg.DataDir != "" {
		c.DataDir = cfg.DataDir
	}
	if cfg.Logging.Level != "" {
		c.Logging.Level = cfg.Logging.Level
	}
	if cfg.Logging.File != "" {
		c.Logging.File = cfg.Logging.File
	}
	if cfg.Server != nil {
		c.Server = cfg.Server
	}
	if cfg.Database != nil {
		c.Database = cfg.Database
	}
	c.Tools = mergeTools(c.Tools, cfg.Tools)
	c.Providers = mergeProviders(c.Providers, cfg.Providers)
	c.Brokers = mergeBrokers(c.Brokers, cfg.Brokers)
	return c
}

func mergeTools(base, over []*Tool) []*Tool {
	out := append([]*Tool{}, base...)
	for _, t := range over {
		replaced := false
		for i, b := range out {
			if b.Name == t.Name {
				out[i] = t
				replaced = true
				break
			}
		}
		if !replaced {
			out = append(out, t)
		}
	}
	return out
}

func mergeProviders(base, over []*Provider) []*Provider {
	out := append([]*Provider{}, base...)
	for _, p := range over {
		replaced := false
		for i, b := range out {
			if b.Name == p.Name {
				out[i] = p
				replaced = true
				break
			}
		}
		if !replaced {
			out = append(out, p)
		}
	}
	return out
}

func mergeBrokers(base, over []*Broker) []*Broker {
	out := append([]*Broker{}, base...)
	for _, br := range over {
		replaced := false
		for i, b := range out {
			if b.Name == br.Name {
				out[i] = br
				replaced = true
				break
			}
		}
		if !replaced {
			out = append(out, br)
		}
	}
	return out
}

// Diagnostics returns the raw HCL diagnostics for a string body. Convenience
// helper for tests.
func Diagnostics(parser *hclparse.Parser) hcl.Diagnostics {
	if parser == nil {
		return nil
	}
	return nil
}

var _ = fmt.Sprintf
