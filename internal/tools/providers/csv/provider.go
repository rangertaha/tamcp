// Package csv provides a "csv" data-source plugin that loads OHLCV bars from
// a local CSV file and exposes a single MCP tool, "csv_bars", that returns
// the requested columns as arrays. Useful for quick local backtests against
// a pre-prepared dataset.
//
// CSV is expected to have at least the columns: sym, open, close, high, low, vol.
package csv

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
)

const Name = "csv"

type pluginCfg struct {
	PricesPath string `hcl:"prices_path,optional"`
	OrdersPath string `hcl:"orders_path,optional"`
}

type plugin struct{}

func (p *plugin) Name() string { return Name }

type input struct {
	Symbol string   `json:"symbol,omitempty" jsonschema:"filter rows by sym (case-insensitive). Empty returns all rows."`
	Fields []string `json:"fields,omitempty" jsonschema:"columns to extract: open, close, high, low, vol. Default = all numeric columns."`
}

type output struct {
	Symbol string               `json:"symbol,omitempty"`
	Count  int                  `json:"count"`
	Series map[string][]float64 `json:"series"`
	Syms   []string             `json:"syms,omitempty"`
}

func (p *plugin) Attach(ctx *tools.Context) error {
	pcfg := pluginCfg{PricesPath: "examples/prices.csv"}
	if ctx.Config != nil {
		if blk := ctx.Config.GetProvider(Name); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &pcfg); diags.HasErrors() {
				return fmt.Errorf("decode csv provider config: %s", diags.Error())
			}
		}
	}

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "csv_bars",
		Description: fmt.Sprintf("Load OHLCV rows from the CSV at %q. Returns each requested numeric column as an array, plus the symbol list.", pcfg.PricesPath),
	}, func(_ context.Context, _ *mcp.CallToolRequest, in input) (*mcp.CallToolResult, output, error) {
		out, err := loadBars(pcfg.PricesPath, in)
		if err != nil {
			return nil, output{}, err
		}
		b, _ := json.Marshal(out)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, out, nil
	})
	return nil
}

func loadBars(path string, in input) (output, error) {
	f, err := os.Open(path)
	if err != nil {
		return output{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return output{}, fmt.Errorf("read header: %w", err)
	}
	col := map[string]int{}
	for i, h := range header {
		col[strings.ToLower(strings.TrimSpace(h))] = i
	}
	want := in.Fields
	if len(want) == 0 {
		for k := range col {
			if k == "sym" {
				continue
			}
			want = append(want, k)
		}
	}

	rows, err := r.ReadAll()
	if err != nil {
		return output{}, fmt.Errorf("read rows: %w", err)
	}

	target := strings.ToUpper(strings.TrimSpace(in.Symbol))
	out := output{Symbol: in.Symbol, Series: map[string][]float64{}}
	syms := map[string]struct{}{}
	for _, row := range rows {
		sym := strings.ToUpper(row[col["sym"]])
		syms[sym] = struct{}{}
		if target != "" && sym != target {
			continue
		}
		out.Count++
		for _, name := range want {
			idx, ok := col[strings.ToLower(name)]
			if !ok {
				continue
			}
			f, _ := strconv.ParseFloat(row[idx], 64)
			out.Series[strings.ToLower(name)] = append(out.Series[strings.ToLower(name)], f)
		}
	}

	if target == "" {
		out.Syms = make([]string, 0, len(syms))
		for s := range syms {
			out.Syms = append(out.Syms, s)
		}
	}
	return out, nil
}

func init() { tools.Register(&plugin{}) }
