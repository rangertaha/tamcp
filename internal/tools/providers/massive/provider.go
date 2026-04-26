// Package massive exposes a polygon.io-shaped market-data provider as MCP
// tools (massive_bars, massive_quote, massive_tickers).
//
// HCL config block:
//
//	provider "massive" {
//	  enabled = true
//	  url     = "https://api.massive.com"
//	  rest {
//	    url    = "https://api.massive.com"
//	    apikey = ""
//	  }
//	}
package massive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
)

const (
	Name           = "massive"
	defaultBaseURL = "https://api.massive.com"
)

// Config is the HCL shape of the `provider "massive" { ... }` block.
type Config struct {
	URL  string      `hcl:"url,optional"`
	Rest *restConfig `hcl:"rest,block"`
}

type restConfig struct {
	URL    string `hcl:"url,optional"`
	APIKey string `hcl:"apikey,optional"`
}

type plugin struct {
	cfg     Config
	baseURL string
	apikey  string
	http    *http.Client
}

func (p *plugin) Name() string { return Name }

func (p *plugin) Attach(ctx *tools.Context) error {
	if ctx.Config != nil {
		if blk := ctx.Config.GetProvider(Name); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &p.cfg); diags.HasErrors() {
				return fmt.Errorf("decode massive config: %s", diags.Error())
			}
		}
	}
	p.baseURL = p.cfg.URL
	if p.cfg.Rest != nil && p.cfg.Rest.URL != "" {
		p.baseURL = p.cfg.Rest.URL
	}
	if strings.TrimSpace(p.baseURL) == "" {
		p.baseURL = defaultBaseURL
	}
	if p.cfg.Rest != nil {
		p.apikey = p.cfg.Rest.APIKey
	}
	p.http = &http.Client{Timeout: 30 * time.Second}

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "massive_bars",
		Description: "Fetch OHLCV aggregates for a symbol/market/interval from massive.com.",
	}, p.handleBars)

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "massive_quote",
		Description: "Fetch the latest quote for a symbol from massive.com.",
	}, p.handleQuote)

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "massive_tickers",
		Description: "Search the massive.com ticker reference catalog.",
	}, p.handleTickers)
	return nil
}

// ── bars ─────────────────────────────────────────────────────

type barsInput struct {
	Symbol   string `json:"symbol" jsonschema:"ticker symbol, e.g. AAPL"`
	Market   string `json:"market,omitempty" jsonschema:"market label (default 'stocks')"`
	Interval string `json:"interval,omitempty" jsonschema:"bar interval, e.g. 1m, 5m, 1h, 1d (default 1d)"`
	From     string `json:"from,omitempty" jsonschema:"start date (RFC3339 or YYYY-MM-DD)"`
	To       string `json:"to,omitempty" jsonschema:"end date (RFC3339 or YYYY-MM-DD)"`
	Limit    int    `json:"limit,omitempty" jsonschema:"max bars (default 120)"`
}

type Bar struct {
	Time   string  `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

type barsOutput struct {
	Symbol   string `json:"symbol"`
	Market   string `json:"market"`
	Interval string `json:"interval"`
	Count    int    `json:"count"`
	Bars     []Bar  `json:"bars"`
}

type rawBars struct {
	Ticker  string `json:"ticker"`
	Results []struct {
		T int64   `json:"t"`
		O float64 `json:"o"`
		H float64 `json:"h"`
		L float64 `json:"l"`
		C float64 `json:"c"`
		V float64 `json:"v"`
	} `json:"results"`
}

func (p *plugin) handleBars(_ context.Context, _ *mcp.CallToolRequest, in barsInput) (*mcp.CallToolResult, barsOutput, error) {
	if strings.TrimSpace(in.Symbol) == "" {
		return nil, barsOutput{}, fmt.Errorf("symbol is required")
	}
	market := in.Market
	if market == "" {
		market = "stocks"
	}
	interval := in.Interval
	if interval == "" {
		interval = "1d"
	}
	limit := in.Limit
	if limit <= 0 {
		limit = 120
	}

	q := url.Values{}
	if in.From != "" {
		q.Set("from", in.From)
	}
	if in.To != "" {
		q.Set("to", in.To)
	}
	q.Set("limit", strconv.Itoa(limit))

	path := fmt.Sprintf("/v2/aggs/%s/%s/%s", strings.ToUpper(in.Symbol), market, interval)

	var raw rawBars
	if err := p.getJSON(path, q, &raw); err != nil {
		return nil, barsOutput{}, err
	}
	out := barsOutput{
		Symbol:   strings.ToUpper(in.Symbol),
		Market:   market,
		Interval: interval,
		Count:    len(raw.Results),
		Bars:     make([]Bar, 0, len(raw.Results)),
	}
	for _, r := range raw.Results {
		out.Bars = append(out.Bars, Bar{
			Time:   time.UnixMilli(r.T).UTC().Format(time.RFC3339),
			Open:   r.O,
			High:   r.H,
			Low:    r.L,
			Close:  r.C,
			Volume: r.V,
		})
	}
	return marshalBars(out)
}

// ── quote ────────────────────────────────────────────────────

type quoteInput struct {
	Symbol string `json:"symbol" jsonschema:"ticker symbol"`
	Market string `json:"market,omitempty" jsonschema:"market label (default 'stocks')"`
}

type Quote struct {
	Symbol    string  `json:"symbol"`
	Bid       float64 `json:"bid"`
	Ask       float64 `json:"ask"`
	Last      float64 `json:"last"`
	Volume    float64 `json:"volume"`
	Timestamp string  `json:"timestamp"`
}

type rawQuote struct {
	Ticker string `json:"ticker"`
	Last   struct {
		Bid       float64 `json:"bid"`
		Ask       float64 `json:"ask"`
		Last      float64 `json:"last"`
		Volume    float64 `json:"volume"`
		Timestamp int64   `json:"t"`
	} `json:"last"`
}

func (p *plugin) handleQuote(_ context.Context, _ *mcp.CallToolRequest, in quoteInput) (*mcp.CallToolResult, Quote, error) {
	if strings.TrimSpace(in.Symbol) == "" {
		return nil, Quote{}, fmt.Errorf("symbol is required")
	}
	market := in.Market
	if market == "" {
		market = "stocks"
	}
	path := fmt.Sprintf("/v2/last/quote/%s/%s", strings.ToUpper(in.Symbol), market)

	var raw rawQuote
	if err := p.getJSON(path, nil, &raw); err != nil {
		return nil, Quote{}, err
	}
	q := Quote{
		Symbol: strings.ToUpper(in.Symbol),
		Bid:    raw.Last.Bid,
		Ask:    raw.Last.Ask,
		Last:   raw.Last.Last,
		Volume: raw.Last.Volume,
	}
	if raw.Last.Timestamp > 0 {
		q.Timestamp = time.UnixMilli(raw.Last.Timestamp).UTC().Format(time.RFC3339)
	}
	return marshalQuote(q)
}

// ── tickers ──────────────────────────────────────────────────

type tickersInput struct {
	Search string `json:"search,omitempty" jsonschema:"substring search across symbol/name"`
	Market string `json:"market,omitempty" jsonschema:"filter by market (stocks, crypto, fx, options)"`
	Active *bool  `json:"active,omitempty" jsonschema:"filter to active tickers (default true)"`
	Limit  int    `json:"limit,omitempty" jsonschema:"max tickers (default 100, max 1000)"`
}

type Ticker struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Market   string `json:"market"`
	Locale   string `json:"locale,omitempty"`
	Exchange string `json:"exchange,omitempty"`
	Type     string `json:"type,omitempty"`
	Active   bool   `json:"active"`
	Currency string `json:"currency,omitempty"`
}

type tickersOutput struct {
	Count   int      `json:"count"`
	Tickers []Ticker `json:"tickers"`
	Next    string   `json:"next,omitempty"`
}

type rawTickers struct {
	Results []struct {
		Ticker          string `json:"ticker"`
		Name            string `json:"name"`
		Market          string `json:"market"`
		Locale          string `json:"locale"`
		PrimaryExchange string `json:"primary_exchange"`
		Type            string `json:"type"`
		Active          bool   `json:"active"`
		Currency        string `json:"currency_name"`
	} `json:"results"`
	Next string `json:"next_url"`
}

func (p *plugin) handleTickers(_ context.Context, _ *mcp.CallToolRequest, in tickersInput) (*mcp.CallToolResult, tickersOutput, error) {
	q := url.Values{}
	limit := in.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	q.Set("limit", strconv.Itoa(limit))
	if in.Search != "" {
		q.Set("search", in.Search)
	}
	if in.Market != "" {
		q.Set("market", in.Market)
	}
	if in.Active != nil {
		q.Set("active", strconv.FormatBool(*in.Active))
	} else {
		q.Set("active", "true")
	}

	var raw rawTickers
	if err := p.getJSON("/v3/reference/tickers", q, &raw); err != nil {
		return nil, tickersOutput{}, err
	}
	out := tickersOutput{Count: len(raw.Results), Next: raw.Next, Tickers: make([]Ticker, 0, len(raw.Results))}
	for _, r := range raw.Results {
		out.Tickers = append(out.Tickers, Ticker{
			Symbol: r.Ticker, Name: r.Name, Market: r.Market, Locale: r.Locale,
			Exchange: r.PrimaryExchange, Type: r.Type, Active: r.Active, Currency: r.Currency,
		})
	}
	return marshalTickers(out)
}

// ── HTTP ─────────────────────────────────────────────────────

func (p *plugin) getJSON(path string, q url.Values, out any) error {
	u, err := url.Parse(strings.TrimRight(p.baseURL, "/") + path)
	if err != nil {
		return err
	}
	if q == nil {
		q = url.Values{}
	}
	if p.apikey != "" {
		q.Set("apikey", p.apikey)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := p.http.Do(req)
	if err != nil {
		return fmt.Errorf("massive: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("massive %s: HTTP %d: %s", path, resp.StatusCode, snippet(body))
	}
	return json.Unmarshal(body, out)
}

func snippet(b []byte) string {
	if len(b) > 240 {
		return string(b[:240]) + "…"
	}
	return string(b)
}

func marshalBars(o barsOutput) (*mcp.CallToolResult, barsOutput, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, o, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func marshalQuote(o Quote) (*mcp.CallToolResult, Quote, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, o, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func marshalTickers(o tickersOutput) (*mcp.CallToolResult, tickersOutput, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, o, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func init() { tools.Register(&plugin{}) }
