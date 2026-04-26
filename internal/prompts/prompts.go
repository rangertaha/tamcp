// Package prompts registers the MCP prompts tamcp exposes. Prompts are
// reusable instruction templates that orchestrate the indicator dispatcher
// and providers into common analytic workflows.
package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Attach registers every prompt with the MCP server.
func Attach(s *mcp.Server) {
	s.AddPrompt(&mcp.Prompt{
		Name:        "analyze_series",
		Description: "Run a battery of trend, momentum, and volatility indicators on a closing-price series and summarize the result.",
		Arguments: []*mcp.PromptArgument{
			{Name: "symbol", Description: "Ticker label for the series", Required: true},
		},
	}, handleAnalyzeSeries)

	s.AddPrompt(&mcp.Prompt{
		Name:        "detect_patterns",
		Description: "Scan an OHLC series for the canonical TA-Lib candlestick patterns and report which bars triggered which patterns.",
		Arguments: []*mcp.PromptArgument{
			{Name: "symbol", Description: "Ticker label", Required: true},
		},
	}, handleDetectPatterns)

	s.AddPrompt(&mcp.Prompt{
		Name:        "trend_summary",
		Description: "Produce a tight trend summary using SMA(20), SMA(50), RSI(14), MACD, and BBANDS(20, 2, 2) on a close series.",
		Arguments: []*mcp.PromptArgument{
			{Name: "symbol", Description: "Ticker label", Required: true},
		},
	}, handleTrendSummary)
}

// ── analyze_series ────────────────────────────────────────────

func handleAnalyzeSeries(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	sym := required(req.Params.Arguments, "symbol")
	if sym == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	text := fmt.Sprintf(`You are a technical analyst. Produce a concise analytic brief for %[1]s.

You have access to the "indicator" MCP tool, which dispatches to all 158 TA-Lib functions.
Call indicator with {"name": "help"} first if you need to see what's available.

Steps:
1. Obtain the close-price series (the user has it; ask if needed).
2. Compute these via indicator:
   - sma (period=20), sma (period=50)
   - ema (period=12), ema (period=26)
   - rsi (period=14)
   - macd (fast_period=12, slow_period=26, signal_period=9)
   - bbands (period=20, nbdevup=2, nbdevdn=2)
   - atr if HLC is available (period=14)
3. Summarize in 6–10 bullets:
   - Trend (SMA20 vs SMA50; EMA crossover state)
   - Momentum (RSI zone; MACD histogram direction)
   - Volatility (BBANDS width; ATR if available)
   - One bullish thesis and one bearish thesis in a single sentence each.

Cite the last observed value for each indicator. Numbers to 2–4 decimals.`,
		strings.ToUpper(sym))

	return promptResult(fmt.Sprintf("Analytic brief for %s", strings.ToUpper(sym)), text), nil
}

// ── detect_patterns ───────────────────────────────────────────

func handleDetectPatterns(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	sym := required(req.Params.Arguments, "symbol")
	if sym == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	text := fmt.Sprintf(`Scan the most recent OHLC bars for %[1]s and report any candlestick patterns that fired.

Use the "indicator" MCP tool to call each pattern. Start with the high-signal set:

  cdldoji, cdlhammer, cdlhangingman, cdlinvertedhammer, cdlshootingstar,
  cdlengulfing, cdlharami, cdlharamicross, cdlmorningstar, cdleveningstar,
  cdl3whitesoldiers, cdl3blackcrows, cdlpiercing, cdldarkcloudcover,
  cdlmarubozu, cdlspinningtop

For each pattern, report any bar where the value is non-zero. Output convention:
  +100 = bullish signal, -100 = bearish signal.

Group findings by bar index (most recent first). For each fired pattern, give:
- the bar's date or index
- the pattern name
- bullish/bearish
- one sentence on what the formation typically implies.`, strings.ToUpper(sym))

	return promptResult(fmt.Sprintf("Pattern scan for %s", strings.ToUpper(sym)), text), nil
}

// ── trend_summary ─────────────────────────────────────────────

func handleTrendSummary(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	sym := required(req.Params.Arguments, "symbol")
	if sym == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	text := fmt.Sprintf(`Produce a tight (3–5 sentence) trend summary for %[1]s.

Compute via the "indicator" tool:
  - sma (period=20), sma (period=50)
  - rsi (period=14)
  - macd (12, 26, 9)
  - bbands (period=20, nbdevup=2, nbdevdn=2)

Answer:
  1. Is price above/below SMA(20) and SMA(50)? Are SMAs aligned (uptrend), inverted (downtrend), or tangled?
  2. Is RSI overbought (>70), oversold (<30), or neutral?
  3. MACD histogram positive or negative; is it expanding or contracting?
  4. BBANDS width: tight (squeeze) or wide?
  5. One-line directional bias.

Quote the last value of each indicator with 2–4 decimals.`, strings.ToUpper(sym))

	return promptResult(fmt.Sprintf("Trend summary for %s", strings.ToUpper(sym)), text), nil
}

// ── helpers ───────────────────────────────────────────────────

func required(args map[string]string, key string) string {
	v, ok := args[key]
	if !ok {
		return ""
	}
	return strings.TrimSpace(v)
}

func promptResult(desc, body string) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: desc,
		Messages: []*mcp.PromptMessage{
			{
				Role:    mcp.Role("user"),
				Content: &mcp.TextContent{Text: body},
			},
		},
	}
}
