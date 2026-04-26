// Package tools is the plugin registry for MCP tools. Each tool subpackage
// registers itself via init() and is wired onto the MCP server at agent
// startup.
package tools

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/config"
	"gorm.io/gorm"
)

// Context is the runtime passed to a plugin when it is attached to the MCP
// server. Plugins add their tools via mcp.AddTool(ctx.Server, ...).
type Context struct {
	Server *mcp.Server
	Config *config.Config
	DB     *gorm.DB
}

// Plugin is the interface every tool plugin implements.
type Plugin interface {
	Name() string
	Attach(ctx *Context) error
}

var (
	registryMu sync.Mutex
	registry   = map[string]Plugin{}
)

// Register adds p to the plugin registry. Call from the plugin's init().
func Register(p Plugin) {
	registryMu.Lock()
	defer registryMu.Unlock()
	name := p.Name()
	if _, dup := registry[name]; dup {
		panic(fmt.Sprintf("tools: duplicate plugin %q", name))
	}
	registry[name] = p
}

// Names returns the names of all registered plugins, sorted.
func Names() []string {
	registryMu.Lock()
	defer registryMu.Unlock()
	out := make([]string, 0, len(registry))
	for n := range registry {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// AttachAll attaches every enabled plugin to the MCP server.
//
// A plugin is enabled if a matching `tool "name" { enabled = true }` or
// `provider "name" { enabled = true }` block is present in the config, or if
// no matching block is declared (default-on).
func AttachAll(ctx *Context) error {
	registryMu.Lock()
	defer registryMu.Unlock()

	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, name := range names {
		if !pluginEnabled(ctx.Config, name) {
			continue
		}
		if err := registry[name].Attach(ctx); err != nil {
			return fmt.Errorf("attach %q: %w", name, err)
		}
	}
	return nil
}

func pluginEnabled(cfg *config.Config, name string) bool {
	if cfg == nil {
		return true
	}
	// Plugin registry keys may carry a suffix to disambiguate when the same
	// label is used in multiple categories (e.g. "polymarket" exists as both a
	// provider and a broker, and registers under "polymarket" + "polymarket_broker").
	// Strip well-known suffixes for the config lookup.
	cfgName := name
	for _, suffix := range []string{"_broker", "_provider"} {
		if strings.HasSuffix(cfgName, suffix) {
			cfgName = strings.TrimSuffix(cfgName, suffix)
			break
		}
	}
	for _, t := range cfg.Tools {
		if t.Name == cfgName {
			return t.Enabled
		}
	}
	if strings.HasSuffix(name, "_broker") || hasBrokerBlock(cfg, cfgName) {
		for _, b := range cfg.Brokers {
			if b.Name == cfgName {
				return b.Enabled
			}
		}
	}
	for _, p := range cfg.Providers {
		if p.Name == cfgName {
			return p.Enabled
		}
	}
	for _, b := range cfg.Brokers {
		if b.Name == cfgName {
			return b.Enabled
		}
	}
	return true
}

func hasBrokerBlock(cfg *config.Config, name string) bool {
	for _, b := range cfg.Brokers {
		if b.Name == name {
			return true
		}
	}
	return false
}
