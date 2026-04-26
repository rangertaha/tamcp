// Package indicator registers a single MCP tool, "indicator", that
// dispatches to any of the technical-analysis functions registered in
// internal/tools/indicators/talib. The catalog covers the full TA-Lib
// surface plus a community set drawn from Pandas TA, sdcoffey/techan,
// cinar/indicator, Ta-Lib-Rust, and Yatala.
//
// Collapsing many narrow tools into one dispatcher keeps the model's tool
// list small (which is the main driver of selection accuracy and token
// overhead). The catalog of supported indicators is discoverable via
// {"name": "help"} and enumerated in this tool's description.
package indicator

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

const Name = "indicator"

type plugin struct{}

func (p *plugin) Name() string { return Name }

type input struct {
	Name string         `json:"name" jsonschema:"indicator name (TA-Lib convention, lower-cased). Pass 'help' to list the catalog or 'help:<name>' for a single indicator's schema."`
	Args map[string]any `json:"args,omitempty" jsonschema:"indicator-specific arguments. Shape depends on 'name'; call with name='help:<name>' to see the schema."`
}

type output struct {
	Result  any    `json:"result"`
	Summary string `json:"summary"`
}

func (p *plugin) Attach(ctx *tools.Context) error {
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        Name,
		Description: buildDescription(),
	}, handle)
	return nil
}

func handle(_ context.Context, _ *mcp.CallToolRequest, in input) (*mcp.CallToolResult, output, error) {
	name := strings.ToLower(strings.TrimSpace(in.Name))
	if name == "" || name == "help" {
		return helpCatalog()
	}
	if strings.HasPrefix(name, "help:") {
		return helpOne(strings.TrimPrefix(name, "help:"))
	}

	entry := talib.Lookup(name)
	if entry == nil {
		hints := talib.SuggestClose(name, 8)
		var msg string
		if len(hints) == 0 {
			msg = fmt.Sprintf("unknown indicator %q; call indicator with name=\"help\" to list supported indicators", name)
		} else {
			msg = fmt.Sprintf("unknown indicator %q; did you mean: %s", name, strings.Join(hints, ", "))
		}
		return errResult(msg)
	}

	args := in.Args
	if args == nil {
		args = map[string]any{}
	}
	out, summary, err := entry.Run(args)
	if err != nil {
		return errResult(err.Error())
	}
	return reply(output{Result: out, Summary: summary})
}

type catalogEntry struct {
	Name        string        `json:"name"`
	Group       string        `json:"group,omitempty"`
	Description string        `json:"description"`
	Params      []talib.Param `json:"params,omitempty"`
}

func helpCatalog() (*mcp.CallToolResult, output, error) {
	all := talib.Catalog()
	cat := make([]catalogEntry, len(all))
	for i, e := range all {
		cat[i] = catalogEntry{Name: e.Name, Group: e.Group, Description: e.Description}
	}
	summary := fmt.Sprintf("%d indicators available — call indicator with name=\"help:<name>\" for the argument schema.", len(cat))
	return reply(output{
		Result:  map[string]any{"count": len(cat), "indicators": cat},
		Summary: summary,
	})
}

func helpOne(name string) (*mcp.CallToolResult, output, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	entry := talib.Lookup(name)
	if entry == nil {
		hints := talib.SuggestClose(name, 8)
		return errResult(fmt.Sprintf("unknown indicator %q; candidates: %s", name, strings.Join(hints, ", ")))
	}
	return reply(output{
		Result: catalogEntry{
			Name:        entry.Name,
			Group:       entry.Group,
			Description: entry.Description,
			Params:      entry.Params,
		},
		Summary: fmt.Sprintf("%s: %s", entry.Name, entry.Description),
	})
}

func buildDescription() string {
	var b strings.Builder
	b.WriteString("Dispatch to any TA-Lib technical indicator. ")
	b.WriteString("Call with {\"name\": \"<indicator>\", \"args\": {...}}. ")
	b.WriteString("Use {\"name\": \"help\"} for the catalog, {\"name\": \"help:<indicator>\"} for a schema.\n\n")
	b.WriteString("Indicators by group:\n")

	groups := map[string][]string{}
	for _, e := range talib.Catalog() {
		g := e.Group
		if g == "" {
			g = "other"
		}
		groups[g] = append(groups[g], e.Name)
	}
	order := make([]string, 0, len(groups))
	for g := range groups {
		order = append(order, g)
	}
	sort.Strings(order)
	for _, g := range order {
		names := groups[g]
		sort.Strings(names)
		b.WriteString("  - ")
		b.WriteString(g)
		b.WriteString(": ")
		b.WriteString(strings.Join(names, ", "))
		b.WriteString("\n")
	}
	return b.String()
}

func reply(o output) (*mcp.CallToolResult, output, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, o, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, o, nil
}

func errResult(msg string) (*mcp.CallToolResult, output, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}, output{Summary: msg}, nil
}

func init() { tools.Register(&plugin{}) }
