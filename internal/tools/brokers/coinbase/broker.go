// Package coinbase is a placeholder broker plugin for Coinbase order
// placement.
//
// HCL config:
//
//	broker "coinbase" {
//	  enabled   = false
//	  url       = "https://api.coinbase.com"
//	  apikey    = ""
//	  apisecret = ""
//	  sandbox   = false
//	}
package coinbase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
)

const Name = "coinbase"

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
		if blk := ctx.Config.GetBroker(Name); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &p.cfg); diags.HasErrors() {
				return fmt.Errorf("decode coinbase broker config: %s", diags.Error())
			}
		}
	}
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "coinbase_place_order",
		Description: "PLACEHOLDER: place a Coinbase order. Not yet implemented.",
	}, p.placeOrder)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "coinbase_cancel_order",
		Description: "PLACEHOLDER: cancel a Coinbase order. Not yet implemented.",
	}, p.cancelOrder)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "coinbase_balance",
		Description: "PLACEHOLDER: fetch Coinbase account balances. Not yet implemented.",
	}, p.balance)
	return nil
}

type placeInput struct {
	ProductID  string  `json:"product_id" jsonschema:"e.g. BTC-USD"`
	Side       string  `json:"side" jsonschema:"buy or sell"`
	Type       string  `json:"type" jsonschema:"market or limit"`
	Size       float64 `json:"size,omitempty" jsonschema:"base-currency quantity"`
	Funds      float64 `json:"funds,omitempty" jsonschema:"quote-currency notional"`
	LimitPrice float64 `json:"limit_price,omitempty"`
}

type cancelInput struct {
	OrderID string `json:"order_id"`
}

type emptyInput struct{}

type placeholder struct {
	Status   string `json:"status"`
	Broker   string `json:"broker"`
	Note     string `json:"note"`
	Sandbox  bool   `json:"sandbox"`
	HasCreds bool   `json:"has_credentials"`
}

func (p *plugin) placeOrder(_ context.Context, _ *mcp.CallToolRequest, _ placeInput) (*mcp.CallToolResult, placeholder, error) {
	return reply(p.stub("Coinbase order placement is a placeholder. Wire it to the Advanced Trade /orders endpoint once API keys + signing are in place."))
}

func (p *plugin) cancelOrder(_ context.Context, _ *mcp.CallToolRequest, _ cancelInput) (*mcp.CallToolResult, placeholder, error) {
	return reply(p.stub("Coinbase order cancellation is a placeholder."))
}

func (p *plugin) balance(_ context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, placeholder, error) {
	return reply(p.stub("Coinbase balance fetching is a placeholder."))
}

func (p *plugin) stub(note string) placeholder {
	return placeholder{
		Status:   "not_implemented",
		Broker:   Name,
		Note:     note,
		Sandbox:  p.cfg.Sandbox,
		HasCreds: p.cfg.APIKey != "" && p.cfg.APISecret != "",
	}
}

func reply(o placeholder) (*mcp.CallToolResult, placeholder, error) {
	b, _ := json.Marshal(o)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func init() { tools.Register(&plugin{}) }
