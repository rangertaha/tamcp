// Package agent wires every registered tool, provider, and prompt onto an
// MCP server and runs it.
package agent

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal"
	"github.com/rangertaha/tamcp/internal/config"
	"github.com/rangertaha/tamcp/internal/db"
	"github.com/rangertaha/tamcp/internal/prompts"
	"github.com/rangertaha/tamcp/internal/tools"
	_ "github.com/rangertaha/tamcp/internal/tools/all"
	"github.com/rangertaha/tamcp/internal/winservice"
	"gorm.io/gorm"
)

// Agent owns the MCP server lifecycle, the GORM database handle, and the
// merged configuration.
type Agent struct {
	cfg    *config.Config
	server *mcp.Server
	gdb    *gorm.DB
}

// New returns a ready-to-run Agent. The variadic cfgs are merged in order.
func New(cfgs ...*config.Config) (*Agent, error) {
	cfg, err := config.New(config.WithConfig(cfgs...))
	if err != nil {
		return nil, err
	}
	return &Agent{cfg: cfg}, nil
}

// Init opens the database and builds the MCP server with every enabled
// plugin and prompt attached.
func (a *Agent) Init() error {
	gdb, err := db.Open(a.cfg)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	a.gdb = gdb

	name := internal.NAME
	if a.cfg != nil && a.cfg.Server != nil && a.cfg.Server.Name != "" {
		name = a.cfg.Server.Name
	}

	a.server = mcp.NewServer(&mcp.Implementation{
		Name:    name,
		Version: internal.Version,
	}, nil)

	if err := tools.AttachAll(&tools.Context{
		Server: a.server,
		Config: a.cfg,
		DB:     a.gdb,
	}); err != nil {
		return err
	}
	prompts.Attach(a.server)
	return nil
}

// Run initializes the agent (if needed) and serves MCP. When launched by the
// Windows Service Control Manager, the server runs under svc.Run so SCM
// stop/shutdown requests cancel the context cleanly.
func (a *Agent) Run() error {
	if a.server == nil {
		if err := a.Init(); err != nil {
			return err
		}
	}

	transport := "stdio"
	if a.cfg != nil && a.cfg.Server != nil && a.cfg.Server.Transport != "" {
		transport = a.cfg.Server.Transport
	}
	if transport != "stdio" {
		return fmt.Errorf("unsupported transport: %q (only stdio supported)", transport)
	}

	if isService, _ := winservice.IsService(); isService {
		return winservice.Run(func(ctx context.Context) error {
			return a.server.Run(ctx, &mcp.StdioTransport{})
		})
	}
	return a.server.Run(context.Background(), &mcp.StdioTransport{})
}
