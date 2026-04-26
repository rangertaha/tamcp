// Package db is the GORM persistence layer used by stateful MCP tools.
package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rangertaha/tamcp/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open returns a GORM handle opened against cfg.Database, running AutoMigrate
// against every model registered in this package.
func Open(cfg *config.Config) (*gorm.DB, error) {
	dsn := resolveDSN(cfg)

	if err := ensureParentDir(dsn); err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := gdb.AutoMigrate(&Ticker{}, &Bar{}, &Order{}); err != nil {
		return nil, fmt.Errorf("migrate schema: %w", err)
	}
	return gdb, nil
}

func resolveDSN(cfg *config.Config) string {
	if env := strings.TrimSpace(os.Getenv("MCPP_DB_DSN")); env != "" {
		return env
	}
	if cfg != nil && cfg.Database != nil && strings.TrimSpace(cfg.Database.DSN) != "" {
		return cfg.Database.DSN
	}
	if cfg != nil && strings.TrimSpace(cfg.DataDir) != "" {
		return filepath.Join(cfg.DataDir, "tamcp.db")
	}
	return "tamcp.db"
}

func ensureParentDir(dsn string) error {
	dir := filepath.Dir(dsn)
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}
