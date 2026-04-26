// Package polymarket is a placeholder broker plugin for Polymarket order
// placement. The HCL config block is read so credentials are available; the
// MCP tools are registered as stubs that return a not-implemented error.
//
// HCL config:
//
//	broker "polymarket" {
//	  enabled   = false
//	  url       = "https://gamma-api.polymarket.com"
//	  apikey    = ""
//	  apisecret = ""
//	}
package polymarket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
)

// Use a distinct registry key from the polymarket *provider* so they coexist.
const Name = "polymarket_broker"

type Config struct {
	URL       string `hcl:"url,optional"`
	APIKey    string `hcl:"apikey,optional"`
	APISecret string `hcl:"apisecret,optional"`
	Sandbox   bool   `hcl:"sandbox,optional"`
}

type plugin struct {
	cfg Config
}

func (p *plugin) Name() string { return Name }

func (p *plugin) Attach(ctx *tools.Context) error {
	if ctx.Config != nil {
		// The HCL block label is "polymarket"; the registry key is "polymarket_broker".
		if blk := ctx.Config.GetBroker("polymarket"); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &p.cfg); diags.HasErrors() {
				return fmt.Errorf("decode polymarket broker config: %s", diags.Error())
			}
		}
	}
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "polymarket_place_order",
		Description: "PLACEHOLDER: place a Polymarket order. Not yet implemented.",
	}, p.placeOrder)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "polymarket_cancel_order",
		Description: "PLACEHOLDER: cancel a Polymarket order. Not yet implemented.",
	}, p.cancelOrder)
	return nil
}

type placeInput struct {
	MarketSlug string  `json:"market_slug" jsonschema:"market slug, e.g. 'will-x-happen'"`
	Side       string  `json:"side" jsonschema:"yes or no"`
	Size       float64 `json:"size" jsonschema:"contract count"`
	LimitPrice float64 `json:"limit_price,omitempty" jsonschema:"limit price (0..1); omit for market"`
}

type cancelInput struct {
	OrderID string `json:"order_id"`
}

type placeholder struct {
	Status   string `json:"status"`
	Provider string `json:"provider"`
	Note     string `json:"note"`
	Sandbox  bool   `json:"sandbox"`
}

func (p *plugin) placeOrder(_ context.Context, _ *mcp.CallToolRequest, _ placeInput) (*mcp.CallToolResult, placeholder, error) {
	return reply(placeholder{
		Status:   "not_implemented",
		Provider: Name,
		Note:     "Polymarket order placement is a placeholder. Wire it to the Gamma API and CLOB once credentials and signing are in place.",
		Sandbox:  p.cfg.Sandbox,
	})
}

func (p *plugin) cancelOrder(_ context.Context, _ *mcp.CallToolRequest, _ cancelInput) (*mcp.CallToolResult, placeholder, error) {
	return reply(placeholder{
		Status:   "not_implemented",
		Provider: Name,
		Note:     "Polymarket order cancellation is a placeholder.",
		Sandbox:  p.cfg.Sandbox,
	})
}

func reply(o placeholder) (*mcp.CallToolResult, placeholder, error) {
	b, _ := json.Marshal(o)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func init() { tools.Register(&plugin{}) }
