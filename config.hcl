# ─────────────────────────────────────────────────────────────
# tamcp – local development configuration
#   Loaded with: tamcp --config ./config.hcl server
# ─────────────────────────────────────────────────────────────

debug = true

datadir = "./data"

logging {
  level = "debug"
  file  = ""   # empty → stderr (stdout is reserved for MCP)
}

server {
  name      = "tamcp"
  transport = "stdio"
}

database {
  driver = "sqlite"
  dsn    = "./data/tamcp.db"
}

# ── Analytic tools ───────────────────────────────────────────

tool "indicator" { enabled = true }

# Language / NLP placeholders.
tool "sentiment" {
  enabled = true
  model   = "default"
}

tool "ner" {
  enabled = true
  model   = "default"
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

# ── Data providers ───────────────────────────────────────────

provider "csv" {
  enabled     = true
  prices_path = "./examples/prices.csv"
  orders_path = "./examples/orders.csv"
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

# ── Brokers (trading execution, placeholders) ────────────────

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
