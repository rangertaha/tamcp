# ─────────────────────────────────────────────────────────────
# tamcp – default configuration
# MCP server for technical analysis (TA-Lib indicators + helpers)
# ─────────────────────────────────────────────────────────────

debug = false

datadir = "/var/lib/tamcp"

logging {
  # Log verbosity: "debug", "info", "warn", "error".
  level = "info"

  # Path to a log file. Leave empty to log to stderr (stdout is reserved for MCP).
  file = "/var/log/tamcp.log"
}

server {
  name      = "tamcp"
  transport = "stdio"
}

database {
  driver = "sqlite"
  dsn    = "/var/lib/tamcp/tamcp.db"
}

# ── Analytic tools ───────────────────────────────────────────
# Disable a tool with: tool "<name>" { enabled = false }

tool "indicator" { enabled = true }

# Language / NLP tools (placeholders until implemented).
tool "sentiment" {
  enabled = true
  model   = "default"
}

tool "ner" {
  enabled = true
  model   = "default"
}

# Charting (gonum/plot).
tool "plots" {
  enabled    = true
  output_dir = "/var/lib/tamcp/plots"
  width      = 800
  height     = 480
  dpi        = 96
}

# ── Data providers ───────────────────────────────────────────

provider "csv" {
  enabled     = true
  prices_path = "examples/prices.csv"
  orders_path = "examples/orders.csv"
}

provider "massive" {
  enabled = false
  url     = "https://api.massive.com"

  rest {
    url    = "https://api.massive.com"
    apikey = ""
  }
}

provider "polymarket" {
  enabled = false
  url     = "https://gamma-api.polymarket.com"
  apikey  = ""
}

provider "kalshi" {
  enabled = false
  url     = "https://api.elections.kalshi.com/trade-api/v2"
  apikey  = ""
}

# ── Brokers (trading execution venues, placeholders) ─────────

broker "polymarket" {
  enabled   = false
  url       = "https://gamma-api.polymarket.com"
  apikey    = ""
  apisecret = ""
}

broker "kalshi" {
  enabled  = false
  url      = "https://api.elections.kalshi.com/trade-api/v2"
  apikey   = ""
  email    = ""
  password = ""
}

broker "coinbase" {
  enabled   = false
  url       = "https://api.coinbase.com"
  apikey    = ""
  apisecret = ""
  sandbox   = false
}
