// Package polymarket exposes Polymarket's Gamma API as MCP tools.
//
// HCL config block:
//
//	provider "polymarket" {
//	  enabled = true
//	  url     = "https://gamma-api.polymarket.com"
//	  apikey  = ""
//	}
package polymarket

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
	Name           = "polymarket"
	defaultBaseURL = "https://gamma-api.polymarket.com"
)

// Config is the HCL shape of the `provider "polymarket" { ... }` block.
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
				return fmt.Errorf("decode polymarket config: %s", diags.Error())
			}
		}
	}
	if strings.TrimSpace(p.cfg.URL) == "" {
		p.cfg.URL = defaultBaseURL
	}
	p.http = &http.Client{Timeout: 30 * time.Second}

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "polymarket_markets",
		Description: "List Polymarket prediction markets via the Gamma API. Returns id, slug, title, category, yes/no prices, volume, and status.",
	}, p.handleMarkets)

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "polymarket_market",
		Description: "Fetch a single Polymarket market by slug or condition ID.",
	}, p.handleMarket)
	return nil
}

// ── markets list ─────────────────────────────────────────────

type marketsInput struct {
	Limit    int    `json:"limit,omitempty" jsonschema:"max markets to return (default 50)"`
	Category string `json:"category,omitempty" jsonschema:"filter by category slug, e.g. politics, sports"`
	Active   *bool  `json:"active,omitempty" jsonschema:"only return active markets (default true)"`
	Closed   *bool  `json:"closed,omitempty" jsonschema:"only return closed markets"`
}

type marketsOutput struct {
	Count   int       `json:"count"`
	Markets []*Market `json:"markets"`
}

// Market is the partial Polymarket market shape we surface. Fields not
// understood by Polymarket's response are dropped.
type Market struct {
	ID          string  `json:"id"`
	ConditionID string  `json:"conditionId,omitempty"`
	Slug        string  `json:"slug"`
	Question    string  `json:"question"`
	Category    string  `json:"category"`
	Active      bool    `json:"active"`
	Closed      bool    `json:"closed"`
	Volume      float64 `json:"volume"`
	Liquidity   float64 `json:"liquidity"`
	YesPrice    float64 `json:"yesPrice,omitempty"`
	NoPrice     float64 `json:"noPrice,omitempty"`
	EndDate     string  `json:"endDate,omitempty"`
	URL         string  `json:"url,omitempty"`
}

type rawMarket struct {
	ID            string `json:"id"`
	ConditionID   string `json:"conditionId"`
	Slug          string `json:"slug"`
	Question      string `json:"question"`
	Category      string `json:"category"`
	Active        bool   `json:"active"`
	Closed        bool   `json:"closed"`
	Volume        string `json:"volume"`
	Liquidity     string `json:"liquidity"`
	OutcomePrices string `json:"outcomePrices"` // JSON-encoded "[\"0.51\",\"0.49\"]"
	EndDate       string `json:"endDate"`
}

func (p *plugin) handleMarkets(_ context.Context, _ *mcp.CallToolRequest, in marketsInput) (*mcp.CallToolResult, marketsOutput, error) {
	q := url.Values{}
	limit := in.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	q.Set("limit", strconv.Itoa(limit))
	if in.Category != "" {
		q.Set("category", in.Category)
	}
	switch {
	case in.Active != nil:
		q.Set("active", strconv.FormatBool(*in.Active))
	case in.Closed == nil || !*in.Closed:
		q.Set("active", "true")
	}
	if in.Closed != nil {
		q.Set("closed", strconv.FormatBool(*in.Closed))
	}

	var raws []rawMarket
	if err := p.getJSON("/markets", q, &raws); err != nil {
		return nil, marketsOutput{}, err
	}
	out := marketsOutput{Count: len(raws), Markets: make([]*Market, 0, len(raws))}
	for _, r := range raws {
		out.Markets = append(out.Markets, normalize(r))
	}
	return jsonResult(out)
}

// ── single market ────────────────────────────────────────────

type marketInput struct {
	Slug        string `json:"slug,omitempty" jsonschema:"market slug (e.g. 'will-x-happen-2026')"`
	ConditionID string `json:"condition_id,omitempty" jsonschema:"condition id (alternative to slug)"`
}

func (p *plugin) handleMarket(_ context.Context, _ *mcp.CallToolRequest, in marketInput) (*mcp.CallToolResult, *Market, error) {
	q := url.Values{}
	switch {
	case in.Slug != "":
		q.Set("slug", in.Slug)
	case in.ConditionID != "":
		q.Set("condition_ids", in.ConditionID)
	default:
		return nil, nil, fmt.Errorf("either slug or condition_id is required")
	}
	q.Set("limit", "1")

	var raws []rawMarket
	if err := p.getJSON("/markets", q, &raws); err != nil {
		return nil, nil, err
	}
	if len(raws) == 0 {
		return nil, nil, fmt.Errorf("market not found")
	}
	m := normalize(raws[0])
	return jsonResultPtr(m)
}

// ── HTTP + helpers ───────────────────────────────────────────

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
		return fmt.Errorf("polymarket: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("polymarket %s: HTTP %d: %s", path, resp.StatusCode, snippet(body))
	}
	return json.Unmarshal(body, out)
}

func normalize(r rawMarket) *Market {
	m := &Market{
		ID:          r.ID,
		ConditionID: r.ConditionID,
		Slug:        r.Slug,
		Question:    r.Question,
		Category:    r.Category,
		Active:      r.Active,
		Closed:      r.Closed,
		EndDate:     r.EndDate,
	}
	m.Volume, _ = strconv.ParseFloat(r.Volume, 64)
	m.Liquidity, _ = strconv.ParseFloat(r.Liquidity, 64)
	if r.OutcomePrices != "" {
		var prices []string
		if json.Unmarshal([]byte(r.OutcomePrices), &prices) == nil {
			if len(prices) >= 1 {
				m.YesPrice, _ = strconv.ParseFloat(prices[0], 64)
			}
			if len(prices) >= 2 {
				m.NoPrice, _ = strconv.ParseFloat(prices[1], 64)
			}
		}
	}
	if r.Slug != "" {
		m.URL = "https://polymarket.com/event/" + r.Slug
	}
	return m
}

func snippet(b []byte) string {
	if len(b) > 240 {
		return string(b[:240]) + "…"
	}
	return string(b)
}

func jsonResult(o marketsOutput) (*mcp.CallToolResult, marketsOutput, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, o, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func jsonResultPtr(m *Market) (*mcp.CallToolResult, *Market, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, m, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, m, nil
}

func init() { tools.Register(&plugin{}) }
