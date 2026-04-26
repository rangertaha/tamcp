// Package all imports every analytic tool plugin for side-effect registration.
//
// Importing this package transitively pulls in:
//   - the indicator dispatcher (one MCP tool routing to all TA-Lib functions),
//   - every per-indicator subpackage that registers an Entry into the talib registry,
//   - every provider plugin (CSV bars, massive, polymarket, kalshi),
//   - every broker placeholder (polymarket, kalshi, coinbase),
//   - every language/NLP plugin (sentiment, ner).
package all

import (
	_ "github.com/rangertaha/tamcp/internal/tools/brokers/all"
	_ "github.com/rangertaha/tamcp/internal/tools/indicator"
	_ "github.com/rangertaha/tamcp/internal/tools/indicators/all"
	_ "github.com/rangertaha/tamcp/internal/tools/languages/all"
	_ "github.com/rangertaha/tamcp/internal/tools/plots"
	_ "github.com/rangertaha/tamcp/internal/tools/providers/all"
)
