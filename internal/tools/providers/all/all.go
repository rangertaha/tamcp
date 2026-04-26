// Package all imports every provider plugin for side-effect registration.
package all

import (
	_ "github.com/rangertaha/tamcp/internal/tools/providers/csv"
	_ "github.com/rangertaha/tamcp/internal/tools/providers/kalshi"
	_ "github.com/rangertaha/tamcp/internal/tools/providers/massive"
	_ "github.com/rangertaha/tamcp/internal/tools/providers/polymarket"
)
