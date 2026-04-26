// Package sentiment is a placeholder NLP plugin for sentiment analysis.
//
// HCL config:
//
//	tool "sentiment" {
//	  enabled = true
//	  model   = "default"  # swap for an actual model id when implementing
//	}
package sentiment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
)

const Name = "sentiment"

type Config struct {
	Model string `hcl:"model,optional"`
}

type plugin struct {
	cfg Config
}

func (p *plugin) Name() string { return Name }

func (p *plugin) Attach(ctx *tools.Context) error {
	if ctx.Config != nil {
		if blk := ctx.Config.GetTool(Name); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &p.cfg); diags.HasErrors() {
				return fmt.Errorf("decode sentiment config: %s", diags.Error())
			}
		}
	}
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        Name,
		Description: "PLACEHOLDER: classify the sentiment of a text as positive, negative, or neutral. Not yet implemented.",
	}, p.handle)
	return nil
}

type input struct {
	Text string `json:"text" jsonschema:"text to classify"`
}

type output struct {
	Status string  `json:"status"`
	Tool   string  `json:"tool"`
	Model  string  `json:"model,omitempty"`
	Note   string  `json:"note"`
	Label  string  `json:"label,omitempty"`
	Score  float64 `json:"score,omitempty"`
}

func (p *plugin) handle(_ context.Context, _ *mcp.CallToolRequest, _ input) (*mcp.CallToolResult, output, error) {
	o := output{
		Status: "not_implemented",
		Tool:   Name,
		Model:  p.cfg.Model,
		Note:   "Sentiment analysis is a placeholder. Wire it to a Go NLP library (e.g. cdipaolo/sentiment) or an upstream classifier endpoint.",
	}
	b, _ := json.Marshal(o)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func init() { tools.Register(&plugin{}) }
