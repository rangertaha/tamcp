// Package kalshi exposes Kalshi's prediction-market API as MCP tools.
//
// HCL config block:
//
//	provider "kalshi" {
//	  enabled = true
//	  url     = "https://api.elections.kalshi.com/trade-api/v2"
//	  apikey  = ""
//	}
package kalshi

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
	Name           = "kalshi"
	defaultBaseURL = "https://api.elections.kalshi.com/trade-api/v2"
)

// Config is the HCL shape of the `provider "kalshi" { ... }` block.
type Config struct {
	URL    string `hcl:"url,optional"`
	APIKey string `hcl:"apikey,optional"`
}

type plugin struct {
	cfg  Config
	http *http.Client
}

func (p *plugin) Name() string { return Name }

func (p *plugin) Attach(ctx *tools.Context) error {
	if ctx.Config != nil {
		if blk := ctx.Config.GetProvider(Name); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &p.cfg); diags.HasErrors() {
				return fmt.Errorf("decode kalshi config: %s", diags.Error())
			}
		}
	}
	if strings.TrimSpace(p.cfg.URL) == "" {
		p.cfg.URL = defaultBaseURL
	}
	p.http = &http.Client{Timeout: 30 * time.Second}

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "kalshi_markets",
		Description: "List Kalshi markets. Returns ticker, title, status, yes/no bid/ask, volume, and close time.",
	}, p.handleMarkets)

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "kalshi_events",
		Description: "List Kalshi events (groups of related markets). Returns event ticker, title, category, and member market tickers.",
	}, p.handleEvents)
	return nil
}

// ── markets ──────────────────────────────────────────────────

type marketsInput struct {
	Limit       int    `json:"limit,omitempty" jsonschema:"max markets (default 100, max 1000)"`
	Status      string `json:"status,omitempty" jsonschema:"open, closed, settled (default open)"`
	EventTicker string `json:"event_ticker,omitempty" jsonschema:"filter to a specific event"`
	Cursor      string `json:"cursor,omitempty" jsonschema:"pagination cursor returned by a prior call"`
}

type marketsOutput struct {
	Count   int             `json:"count"`
	Markets []*KalshiMarket `json:"markets"`
	Cursor  string          `json:"cursor,omitempty"`
}

// KalshiMarket is the partial Kalshi market shape we surface.
type KalshiMarket struct {
	Ticker         string  `json:"ticker"`
	EventTicker    string  `json:"event_ticker"`
	Title          string  `json:"title"`
	Status         string  `json:"status"`
	YesBid         int     `json:"yes_bid"` // cents
	YesAsk         int     `json:"yes_ask"`
	NoBid          int     `json:"no_bid"`
	NoAsk          int     `json:"no_ask"`
	LastPrice      int     `json:"last_price"`
	Volume         int64   `json:"volume"`
	OpenInterest   int64   `json:"open_interest"`
	CloseTime      string  `json:"close_time,omitempty"`
	YesProbability float64 `json:"yes_probability,omitempty"`
}

type marketsResp struct {
	Markets []rawMarket `json:"markets"`
	Cursor  string      `json:"cursor"`
}

type rawMarket struct {
	Ticker       string `json:"ticker"`
	EventTicker  string `json:"event_ticker"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	YesBid       int    `json:"yes_bid"`
	YesAsk       int    `json:"yes_ask"`
	NoBid        int    `json:"no_bid"`
	NoAsk        int    `json:"no_ask"`
	LastPrice    int    `json:"last_price"`
	Volume       int64  `json:"volume"`
	OpenInterest int64  `json:"open_interest"`
	CloseTime    string `json:"close_time"`
}

func (p *plugin) handleMarkets(_ context.Context, _ *mcp.CallToolRequest, in marketsInput) (*mcp.CallToolResult, marketsOutput, error) {
	q := url.Values{}
	limit := in.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	q.Set("limit", strconv.Itoa(limit))
	if in.Status != "" {
		q.Set("status", in.Status)
	} else {
		q.Set("status", "open")
	}
	if in.EventTicker != "" {
		q.Set("event_ticker", in.EventTicker)
	}
	if in.Cursor != "" {
		q.Set("cursor", in.Cursor)
	}

	var resp marketsResp
	if err := p.getJSON("/markets", q, &resp); err != nil {
		return nil, marketsOutput{}, err
	}
	out := marketsOutput{Count: len(resp.Markets), Cursor: resp.Cursor, Markets: make([]*KalshiMarket, 0, len(resp.Markets))}
	for _, r := range resp.Markets {
		out.Markets = append(out.Markets, normalize(r))
	}
	return marshalResult(out)
}

// ── events ───────────────────────────────────────────────────

type eventsInput struct {
	Limit  int    `json:"limit,omitempty" jsonschema:"max events (default 100, max 200)"`
	Status string `json:"status,omitempty" jsonschema:"open, closed, settled (default open)"`
	Cursor string `json:"cursor,omitempty"`
}

type eventsOutput struct {
	Count  int            `json:"count"`
	Events []*KalshiEvent `json:"events"`
	Cursor string         `json:"cursor,omitempty"`
}

type KalshiEvent struct {
	EventTicker string   `json:"event_ticker"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Status      string   `json:"status"`
	Markets     []string `json:"markets,omitempty"`
}

type eventsResp struct {
	Events []rawEvent `json:"events"`
	Cursor string     `json:"cursor"`
}

type rawEvent struct {
	EventTicker string      `json:"event_ticker"`
	Title       string      `json:"title"`
	Category    string      `json:"category"`
	Status      string      `json:"status"`
	Markets     []rawMarket `json:"markets"`
}

func (p *plugin) handleEvents(_ context.Context, _ *mcp.CallToolRequest, in eventsInput) (*mcp.CallToolResult, eventsOutput, error) {
	q := url.Values{}
	limit := in.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
	}
	q.Set("limit", strconv.Itoa(limit))
	if in.Status != "" {
		q.Set("status", in.Status)
	} else {
		q.Set("status", "open")
	}
	if in.Cursor != "" {
		q.Set("cursor", in.Cursor)
	}

	var resp eventsResp
	if err := p.getJSON("/events", q, &resp); err != nil {
		return nil, eventsOutput{}, err
	}
	out := eventsOutput{Count: len(resp.Events), Cursor: resp.Cursor, Events: make([]*KalshiEvent, 0, len(resp.Events))}
	for _, r := range resp.Events {
		ev := &KalshiEvent{
			EventTicker: r.EventTicker,
			Title:       r.Title,
			Category:    r.Category,
			Status:      r.Status,
		}
		for _, m := range r.Markets {
			ev.Markets = append(ev.Markets, m.Ticker)
		}
		out.Events = append(out.Events, ev)
	}
	return marshalResultEvents(out)
}

// ── HTTP ─────────────────────────────────────────────────────

func (p *plugin) getJSON(path string, q url.Values, out any) error {
	u, err := url.Parse(strings.TrimRight(p.cfg.URL, "/") + path)
	if err != nil {
		return err
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	if p.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
	}
	resp, err := p.http.Do(req)
	if err != nil {
		return fmt.Errorf("kalshi: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("kalshi %s: HTTP %d: %s", path, resp.StatusCode, snippet(body))
	}
	return json.Unmarshal(body, out)
}

func normalize(r rawMarket) *KalshiMarket {
	m := &KalshiMarket{
		Ticker:       r.Ticker,
		EventTicker:  r.EventTicker,
		Title:        r.Title,
		Status:       r.Status,
		YesBid:       r.YesBid,
		YesAsk:       r.YesAsk,
		NoBid:        r.NoBid,
		NoAsk:        r.NoAsk,
		LastPrice:    r.LastPrice,
		Volume:       r.Volume,
		OpenInterest: r.OpenInterest,
		CloseTime:    r.CloseTime,
	}
	// Mid yes-price as a probability in [0,1].
	if r.YesBid > 0 || r.YesAsk > 0 {
		m.YesProbability = float64(r.YesBid+r.YesAsk) / 200.0
	}
	return m
}

func snippet(b []byte) string {
	if len(b) > 240 {
		return string(b[:240]) + "…"
	}
	return string(b)
}

func marshalResult(o marketsOutput) (*mcp.CallToolResult, marketsOutput, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, o, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func marshalResultEvents(o eventsOutput) (*mcp.CallToolResult, eventsOutput, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, o, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func init() { tools.Register(&plugin{}) }
