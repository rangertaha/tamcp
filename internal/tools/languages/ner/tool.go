// Package ner is a placeholder NLP plugin for Named-Entity Recognition.
//
// HCL config:
//
//	tool "ner" {
//	  enabled = true
//	  model   = "default"
//	}
package ner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
)

const Name = "ner"

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
				return fmt.Errorf("decode ner config: %s", diags.Error())
			}
		}
	}
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        Name,
		Description: "PLACEHOLDER: extract named entities (persons, organizations, locations, money, dates) from text. Not yet implemented.",
	}, p.handle)
	return nil
}

type input struct {
	Text string `json:"text" jsonschema:"text to extract entities from"`
}

type Entity struct {
	Text  string `json:"text"`
	Label string `json:"label"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type output struct {
	Status   string   `json:"status"`
	Tool     string   `json:"tool"`
	Model    string   `json:"model,omitempty"`
	Note     string   `json:"note"`
	Entities []Entity `json:"entities,omitempty"`
}

func (p *plugin) handle(_ context.Context, _ *mcp.CallToolRequest, _ input) (*mcp.CallToolResult, output, error) {
	o := output{
		Status: "not_implemented",
		Tool:   Name,
		Model:  p.cfg.Model,
		Note:   "Named-entity recognition is a placeholder. Wire it to a Go NLP library or an upstream NER endpoint.",
	}
	b, _ := json.Marshal(o)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func init() { tools.Register(&plugin{}) }
