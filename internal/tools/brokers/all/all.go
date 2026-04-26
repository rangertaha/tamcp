// Package all imports every broker plugin for side-effect registration.
package all

import (
	_ "github.com/rangertaha/tamcp/internal/tools/brokers/coinbase"
	_ "github.com/rangertaha/tamcp/internal/tools/brokers/kalshi"
	_ "github.com/rangertaha/tamcp/internal/tools/brokers/polymarket"
)
