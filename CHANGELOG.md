# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Renamed the project, binary, module, and on-disk paths from `mcpp` to `tamcp`.
  Module path is now `github.com/rangertaha/tamcp`; default config locations
  moved to `/etc/tamcp/`, `~/.config/tamcp/`, default DB to `./data/tamcp.db`,
  and the systemd unit to `tamcp.service`.

## [0.4.0] - 2026-04-25

### Added
- MCP server (`tamcp`) speaking JSON-RPC over stdio per the Model Context Protocol.
- `indicator` dispatcher tool exposing 330 indicators behind a single MCP tool, with
  `help` and `help:<name>` discovery.
  - 162 TA-Lib v0.4 functions (overlap, momentum, volume, volatility, statistic,
    cycle, candlestick patterns, math).
  - 168 community indicators sourced from Pandas TA, sdcoffey/techan,
    cinar/indicator, Ta-Lib-Rust, and Yatala (supertrend, ichimoku, vwap, kdj, wt,
    squeeze/squeeze_pro, fisher, vortex, chande_kroll, chandelier_exit,
    alligator/gator, frama, vidya, hwma/hwc, tdi, qqe, ssf, gaussian, cyber_cycle,
    decycler, reflex/trendflex, pivot variants, kc/kcb, donchian, bbp/bbw,
    bb_squeeze, crsi, tsi/kst/coppock, cmf/mfi_signal/obv_signal, and more).
- Plot tools backed by gonum/plot, writing PNGs to `data/plots/` and returning a
  base64 inline copy: `plot_line`, `plot_scatter`, `plot_histogram`, `plot_bar`,
  `plot_terminal`, `plot_csv_line`, `plot_csv_scatter`, `plot_csv_histogram`.
- Optional `auto_open` for plot tools (`xdg-open` / `open` / `rundll32`).
- Data providers: `csv_bars`, `massive_bars` / `massive_quote` / `massive_tickers`,
  `polymarket_markets` / `polymarket_market`, `kalshi_markets` / `kalshi_events`.
- Broker stubs: `coinbase_*`, `polymarket_*`, `kalshi_*` (disabled by default).
- NLP placeholders: `sentiment`, `ner`.
- HCL configuration with global / user / project merge and `-c <path>` override.
- SQLite persistence via GORM (`./data/tamcp.db`).
- Windows Service Control Manager integration (`internal/winservice`).
- Makefile targets: `run`, `build`, `init`, `server`, `test`, `fmt`, `vet`,
  `tidy`, `clean`, `bump` (semver tag bumper).
- Build-time version metadata (`Version`, `Commit`, `BuildDate`) wired via ldflags.

[Unreleased]: https://github.com/rangertaha/tamcp/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/rangertaha/tamcp/releases/tag/v0.4.0
