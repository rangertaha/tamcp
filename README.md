# TAMCP

Technical Analysis Model Context Protocal




`tamcp` is a Go-based MCP server for technical-analysis indicators, charting, market-data providers, and broker integrations. It speaks JSON-RPC over stdio per the Model Context Protocol and is consumable by Claude Code, Claude Desktop, or any MCP-compatible client.

- **330 indicators** in a single dispatcher tool
- **Plot tools** (line, scatter, histogram, bar, terminal sparkline, CSV-sourced variants) that save PNGs and optionally open them in the OS default viewer
- **Data providers**: CSV, Massive, Polymarket, Kalshi
- **Broker stubs**: Coinbase, Polymarket, Kalshi
- **NLP placeholders**: sentiment, NER

## Build and run

```sh
make build                            # → ./bin/tamcp
./bin/tamcp --config ./config.hcl server
```

Or register with Claude Code:

```sh
claude mcp add tamcp /absolute/path/to/bin/tamcp server
```

## Tools

The server exposes one MCP tool per concern, plus a single dispatcher tool for the technical indicator catalog.

### `indicator` — dispatcher

330 indicators behind one tool. Call shape:

```json
{ "name": "rsi", "args": { "values": [...], "period": 14 } }
```

Discoverability:

- `{"name": "help"}` — list every supported indicator
- `{"name": "help:<name>"}` — schema and description for one indicator

The catalog is composed of:

- **162 TA-Lib functions** — full TA-Lib v0.4 surface (overlap, momentum, volume, volatility, statistic, cycle, candlestick patterns, math)
- **168 community indicators** sourced from Pandas TA, sdcoffey/techan, cinar/indicator, Ta-Lib-Rust, and Yatala. Highlights: `supertrend`, `ichimoku`, `vwap`, `vwap_anchored`, `kdj`, `wt` (Wave Trend), `squeeze`/`squeeze_pro`, `fisher`, `vortex`, `chande_kroll`, `chandelier_exit`, `alligator`/`gator`, `frama`, `vidya`, `hwma`/`hwc`, `tdi`, `qqe`, `ssf`, `gaussian`, `cyber_cycle`, `decycler`, `reflex`/`trendflex`, `pivot_cpr`, `camarilla`/`woodie`/`fib_pivots`/`demark_pivots`, `kc`/`kcb`, `donchian`/`donchian_pct`, `bbp`/`bbw`/`bb_squeeze`, `crsi`, `tsi`/`kst`/`coppock`, `cmf`/`mfi_signal`/`obv_signal`/etc.

### Plot tools

Each plot writes a PNG to `data/plots/` (configurable) and returns its path plus a base64 inline copy. When `auto_open = true` (the default), the saved PNG is also opened in the OS default image viewer (`xdg-open` / `open` / `rundll32`).

| Tool | Input shape | Notes |
|---|---|---|
| `plot_line` | `series[]` of `{name, y[], x?[]}` | Multi-series line chart |
| `plot_scatter` | same | XY scatter |
| `plot_histogram` | `values[]`, `bins?` | Frequency distribution |
| `plot_bar` | `labels[]`, `values[]` | Vertical bar chart |
| `plot_terminal` | `series[]` | ANSI/Unicode sparkline (text only, no PNG) |
| `plot_csv_line` | `path`, `y_cols[]`, `symbol?`, `x_col?` | Reads a CSV, plots columns directly |
| `plot_csv_scatter` | same | XY scatter from CSV |
| `plot_csv_histogram` | `path`, `column`, `symbol?`, `bins?` | Histogram from one CSV column |

### Data providers

| Tool | Source | Description |
|---|---|---|
| `csv_bars` | local CSV | OHLCV columns from the configured `provider "csv" { prices_path = ... }` |
| `massive_bars` / `massive_quote` / `massive_tickers` | Massive REST | Bars / latest quote / ticker list |
| `polymarket_markets` / `polymarket_market` | Polymarket Gamma | Prediction-market data |
| `kalshi_markets` / `kalshi_events` | Kalshi | Prediction-market data |

### Broker tools

`coinbase_*`, `polymarket_*`, `kalshi_*` — order placement, balances, market lookups. All disabled by default in `config.hcl`.

### NLP

`sentiment`, `ner` — placeholders pending a model-backed implementation.

## Configuration

Global / user / project HCL files are merged at startup. Local override goes via `-c <path>`.

```hcl
debug   = true
datadir = "./data"

logging {
  level = "debug"
  file  = ""    # empty → stderr (stdout is reserved for MCP)
}

server {
  name      = "tamcp"
  transport = "stdio"
}

database {
  driver = "sqlite"
  dsn    = "./data/tamcp.db"
}

# Charting (gonum/plot). Outputs PNGs to output_dir.
# auto_open pops the saved PNG in the OS default image viewer
# after each plot tool call. Set false on headless servers.
tool "plots" {
  enabled    = true
  output_dir = "./data/plots"
  width      = 800
  height     = 480
  dpi        = 96
  auto_open  = true
}

provider "csv" {
  enabled     = true
  prices_path = "./examples/prices.csv"
  orders_path = "./examples/orders.csv"
}

# Provider/broker blocks (massive, polymarket, kalshi, coinbase) are
# disabled until you supply credentials; see config.hcl for the full surface.
```

## Project layout

```
cmd/tamcp/                 CLI entry point
internal/
  agent/                  MCP server bootstrap
  config/                 HCL config loader
  db/                     SQLite via GORM
  prompts/                System prompts
  tools/                  MCP tool plugins
    indicator/            "indicator" dispatcher tool
    indicators/           One package per indicator (registers with talib dispatcher)
      talib/              Math implementations + dispatcher registry
      all/                Side-effect imports of every indicator subpackage
    plots/                Plot tools (line, scatter, histogram, bar, terminal, csv_*)
    providers/            csv, massive, polymarket, kalshi
    brokers/              coinbase, polymarket, kalshi
    languages/            sentiment, ner
  winservice/             Windows Service Control Manager integration
examples/                 Sample data (prices.csv, orders.csv, pivot.csv)
config.hcl                Local development config
```

## Adding an indicator

1. Add the math function to `internal/tools/indicators/talib/community.go` (or create a new file in that package).
2. Create `internal/tools/indicators/<name>/tool.go` with an `init()` that calls `talib.RegisterEntry(...)`. Use the helpers in `talib/runners.go` (`ParamsRealPeriod`, `RunHLCPeriod`, etc.) for the common signatures.
3. Add the import to `internal/tools/indicators/all/all.go`.

The dispatcher tool's description and `help` catalog update automatically.

## Library inspirations

The community indicator set draws on:

- [Pandas TA](https://github.com/twopirllc/pandas-ta)
- [sdcoffey/techan](https://github.com/sdcoffey/techan)
- [cinar/indicator](https://github.com/cinar/indicator)
- [Ta-Lib-Rust](https://github.com/greyblake/ta-rs)
- Yatala
