// Package kalshi is a placeholder broker plugin for Kalshi order placement.
//
// HCL config:
//
//	broker "kalshi" {
//	  enabled = false
//	  url     = "https://api.elections.kalshi.com/trade-api/v2"
//	  apikey  = ""
//	  email   = ""
//	  password = ""
//	}
package kalshi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
)

// Use a distinct registry key from the kalshi *provider* so they coexist.
const Name = "kalshi_broker"

type Config struct {
	URL      string `hcl:"url,optional"`
	APIKey   string `hcl:"apikey,optional"`
	Email    string `hcl:"email,optional"`
	Password string `hcl:"password,optional"`
}

type plugin struct {
	cfg Config
}

func (p *plugin) Name() string { return Name }

func (p *plugin) Attach(ctx *tools.Context) error {
	if ctx.Config != nil {
		// Read the `broker "kalshi"` block (the broker label is "kalshi", not the registry name).
		if blk := ctx.Config.GetBroker("kalshi"); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &p.cfg); diags.HasErrors() {
				return fmt.Errorf("decode kalshi broker config: %s", diags.Error())
			}
		}
	}
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "kalshi_place_order",
		Description: "PLACEHOLDER: place a Kalshi order. Not yet implemented.",
	}, p.placeOrder)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "kalshi_cancel_order",
		Description: "PLACEHOLDER: cancel a Kalshi order. Not yet implemented.",
	}, p.cancelOrder)
	return nil
}

type placeInput struct {
	MarketTicker string `json:"market_ticker"`
	Side         string `json:"side" jsonschema:"yes or no"`
	Action       string `json:"action" jsonschema:"buy or sell"`
	Count        int    `json:"count"`
	YesPrice     int    `json:"yes_price,omitempty" jsonschema:"price in cents (1-99) for yes; omit for market"`
}

type cancelInput struct {
	OrderID string `json:"order_id"`
}

type placeholder struct {
	Status   string `json:"status"`
	Broker   string `json:"broker"`
	Note     string `json:"note"`
	HasCreds bool   `json:"has_credentials"`
}

func (p *plugin) placeOrder(_ context.Context, _ *mcp.CallToolRequest, _ placeInput) (*mcp.CallToolResult, placeholder, error) {
	return reply(placeholder{
		Status:   "not_implemented",
		Broker:   "kalshi",
		Note:     "Kalshi order placement is a placeholder. Wire it to /trade-api/v2/portfolio/orders once auth (login + bearer) is implemented.",
		HasCreds: p.cfg.APIKey != "" || (p.cfg.Email != "" && p.cfg.Password != ""),
	})
}

func (p *plugin) cancelOrder(_ context.Context, _ *mcp.CallToolRequest, _ cancelInput) (*mcp.CallToolResult, placeholder, error) {
	return reply(placeholder{
		Status: "not_implemented",
		Broker: "kalshi",
		Note:   "Kalshi order cancellation is a placeholder.",
	})
}

func reply(o placeholder) (*mcp.CallToolResult, placeholder, error) {
	b, _ := json.Marshal(o)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func init() { tools.Register(&plugin{}) }
