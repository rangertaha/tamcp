// Package plots renders data visualizations via gonum/plot and returns each
// chart as a PNG (file on disk + base64 inline so MCP clients can display it).
//
// Tools registered:
//   - plot_line          — multiple line series on one chart
//   - plot_scatter       — XY scatter
//   - plot_histogram     — frequency distribution of a series
//   - plot_bar           — vertical bar chart
//   - plot_csv_line      — line chart sourced from a CSV file
//   - plot_csv_scatter   — scatter sourced from a CSV file
//   - plot_csv_histogram — histogram sourced from one CSV column
//
// HCL config:
//
//	tool "plots" {
//	  enabled    = true
//	  output_dir = "./data/plots"   # PNGs written here
//	  width      = 800              # default canvas width  (pixels)
//	  height     = 480              # default canvas height (pixels)
//	}
package plots

import (
	"context"
	"encoding/base64"
	stdcsv "encoding/csv"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/tamcp/internal/tools"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const Name = "plots"

type Config struct {
	OutputDir string  `hcl:"output_dir,optional"`
	Width     int     `hcl:"width,optional"`
	Height    int     `hcl:"height,optional"`
	DPI       float64 `hcl:"dpi,optional"`
	// AutoOpen pops the saved PNG in the OS default image viewer after
	// each plot tool call. Set false on headless servers. Defaults to true.
	AutoOpen *bool `hcl:"auto_open,optional"`
}

type plugin struct {
	cfg Config
}

func (p *plugin) Name() string { return Name }

func (p *plugin) Attach(ctx *tools.Context) error {
	if ctx.Config != nil {
		if blk := ctx.Config.GetTool(Name); blk != nil && blk.Body != nil {
			if diags := gohcl.DecodeBody(blk.Body, nil, &p.cfg); diags.HasErrors() {
				return fmt.Errorf("decode plots config: %s", diags.Error())
			}
		}
	}
	if strings.TrimSpace(p.cfg.OutputDir) == "" {
		p.cfg.OutputDir = "./data/plots"
	}
	if p.cfg.Width <= 0 {
		p.cfg.Width = 800
	}
	if p.cfg.Height <= 0 {
		p.cfg.Height = 480
	}
	if p.cfg.DPI <= 0 {
		p.cfg.DPI = 96
	}
	if p.cfg.AutoOpen == nil {
		t := true
		p.cfg.AutoOpen = &t
	}

	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_line",
		Description: "Render one or more line series as a PNG. Returns the file path and the image inline (base64).",
	}, p.handleLine)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_scatter",
		Description: "Render an XY scatter plot.",
	}, p.handleScatter)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_histogram",
		Description: "Render a histogram (frequency distribution) of a numeric series.",
	}, p.handleHistogram)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_bar",
		Description: "Render a vertical bar chart from labeled values.",
	}, p.handleBar)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_terminal",
		Description: "Render line series as an ANSI/Unicode sparkline chart that prints in a terminal. Returns the chart as plain text content (no image).",
	}, p.handleTerminal)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_csv_line",
		Description: "Read a CSV file, optionally filter by a symbol column, and plot one or more numeric columns as a line chart. PNG saved to output_dir; opens in the default viewer when auto_open is true.",
	}, p.handleCSVLine)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_csv_scatter",
		Description: "Read a CSV file, optionally filter by a symbol column, and plot one or more numeric columns as an XY scatter.",
	}, p.handleCSVScatter)
	mcp.AddTool(ctx.Server, &mcp.Tool{
		Name:        "plot_csv_histogram",
		Description: "Read a CSV file, optionally filter by a symbol column, and plot a histogram of one numeric column.",
	}, p.handleCSVHistogram)
	return nil
}

// ── terminal ─────────────────────────────────────────────────

type terminalInput struct {
	Title   string   `json:"title,omitempty"`
	Width   int      `json:"width,omitempty" jsonschema:"chart width in cells (default 80)"`
	Height  int      `json:"height,omitempty" jsonschema:"chart height in rows (default 16)"`
	Series  []Series `json:"series" jsonschema:"one or more series; if multiple, each is plotted in a different color"`
	Colored bool     `json:"colored,omitempty" jsonschema:"emit ANSI color codes (default true)"`
}

type terminalOutput struct {
	Chart string `json:"chart"`
	Width int    `json:"width"`
	Rows  int    `json:"rows"`
}

func (p *plugin) handleTerminal(_ context.Context, _ *mcp.CallToolRequest, in terminalInput) (*mcp.CallToolResult, terminalOutput, error) {
	if len(in.Series) == 0 {
		return nil, terminalOutput{}, fmt.Errorf("at least one series is required")
	}
	width := in.Width
	if width <= 0 {
		width = 80
	}
	height := in.Height
	if height <= 0 {
		height = 16
	}

	data := make([][]float64, len(in.Series))
	for i, s := range in.Series {
		data[i] = s.Y
	}
	opts := []asciigraph.Option{
		asciigraph.Width(width),
		asciigraph.Height(height),
	}
	if in.Title != "" {
		opts = append(opts, asciigraph.Caption(in.Title))
	}
	// Colour every series distinctly when more than one is given.
	if len(in.Series) > 1 || in.Colored {
		palette := []asciigraph.AnsiColor{
			asciigraph.Red, asciigraph.Green, asciigraph.Blue,
			asciigraph.Yellow, asciigraph.Magenta, asciigraph.Cyan,
			asciigraph.White,
		}
		colors := make([]asciigraph.AnsiColor, len(in.Series))
		for i := range in.Series {
			colors[i] = palette[i%len(palette)]
		}
		opts = append(opts, asciigraph.SeriesColors(colors...))
	}

	chart := asciigraph.PlotMany(data, opts...)
	// Prepend a series legend since asciigraph's caption doesn't carry names.
	if len(in.Series) > 1 {
		var b strings.Builder
		for i, s := range in.Series {
			fmt.Fprintf(&b, "  [%d] %s   n=%d\n", i+1, s.Name, len(s.Y))
		}
		chart = b.String() + chart
	}
	out := terminalOutput{Chart: chart, Width: width, Rows: height}
	body, _ := json.Marshal(map[string]any{"width": width, "rows": height})
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(body)},
			&mcp.TextContent{Text: chart},
		},
	}, out, nil
}

// ── line ─────────────────────────────────────────────────────

type Series struct {
	Name string    `json:"name" jsonschema:"series label (used in legend)"`
	Y    []float64 `json:"y" jsonschema:"y values"`
	X    []float64 `json:"x,omitempty" jsonschema:"optional x values; defaults to 0..len(y)-1"`
}

type lineInput struct {
	Title    string   `json:"title,omitempty"`
	XLabel   string   `json:"x_label,omitempty"`
	YLabel   string   `json:"y_label,omitempty"`
	Filename string   `json:"filename,omitempty" jsonschema:"output filename (without dir; defaults to a slug of title + timestamp)"`
	Series   []Series `json:"series" jsonschema:"one or more series to plot"`
}

func (p *plugin) handleLine(_ context.Context, _ *mcp.CallToolRequest, in lineInput) (*mcp.CallToolResult, output, error) {
	if len(in.Series) == 0 {
		return nil, output{}, fmt.Errorf("at least one series is required")
	}
	pl := plot.New()
	applyTitles(pl, in.Title, in.XLabel, in.YLabel)

	args := make([]any, 0, 2*len(in.Series))
	for _, s := range in.Series {
		args = append(args, s.Name, toXYs(s))
	}
	if err := plotutil.AddLines(pl, args...); err != nil {
		return nil, output{}, fmt.Errorf("add lines: %w", err)
	}
	return p.save(pl, in.Filename, in.Title, "line")
}

// ── scatter ──────────────────────────────────────────────────

type scatterInput struct {
	Title    string   `json:"title,omitempty"`
	XLabel   string   `json:"x_label,omitempty"`
	YLabel   string   `json:"y_label,omitempty"`
	Filename string   `json:"filename,omitempty"`
	Series   []Series `json:"series" jsonschema:"one or more series to plot"`
}

func (p *plugin) handleScatter(_ context.Context, _ *mcp.CallToolRequest, in scatterInput) (*mcp.CallToolResult, output, error) {
	if len(in.Series) == 0 {
		return nil, output{}, fmt.Errorf("at least one series is required")
	}
	pl := plot.New()
	applyTitles(pl, in.Title, in.XLabel, in.YLabel)

	args := make([]any, 0, 2*len(in.Series))
	for _, s := range in.Series {
		args = append(args, s.Name, toXYs(s))
	}
	if err := plotutil.AddScatters(pl, args...); err != nil {
		return nil, output{}, fmt.Errorf("add scatters: %w", err)
	}
	return p.save(pl, in.Filename, in.Title, "scatter")
}

// ── histogram ────────────────────────────────────────────────

type histogramInput struct {
	Title    string    `json:"title,omitempty"`
	XLabel   string    `json:"x_label,omitempty"`
	YLabel   string    `json:"y_label,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Values   []float64 `json:"values" jsonschema:"observations to bin"`
	Bins     int       `json:"bins,omitempty" jsonschema:"number of bins (default 30)"`
}

func (p *plugin) handleHistogram(_ context.Context, _ *mcp.CallToolRequest, in histogramInput) (*mcp.CallToolResult, output, error) {
	if len(in.Values) == 0 {
		return nil, output{}, fmt.Errorf("values is required")
	}
	pl := plot.New()
	applyTitles(pl, in.Title, in.XLabel, in.YLabel)
	if in.YLabel == "" {
		pl.Y.Label.Text = "frequency"
	}
	bins := in.Bins
	if bins <= 0 {
		bins = 30
	}
	v := plotter.Values(in.Values)
	h, err := plotter.NewHist(v, bins)
	if err != nil {
		return nil, output{}, fmt.Errorf("new histogram: %w", err)
	}
	h.FillColor = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	pl.Add(h)
	return p.save(pl, in.Filename, in.Title, "histogram")
}

// ── bar ──────────────────────────────────────────────────────

type barInput struct {
	Title    string    `json:"title,omitempty"`
	XLabel   string    `json:"x_label,omitempty"`
	YLabel   string    `json:"y_label,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Labels   []string  `json:"labels" jsonschema:"x-axis category labels"`
	Values   []float64 `json:"values" jsonschema:"bar heights; same length as labels"`
}

func (p *plugin) handleBar(_ context.Context, _ *mcp.CallToolRequest, in barInput) (*mcp.CallToolResult, output, error) {
	if len(in.Labels) == 0 || len(in.Values) == 0 || len(in.Labels) != len(in.Values) {
		return nil, output{}, fmt.Errorf("labels and values must be non-empty and the same length")
	}
	pl := plot.New()
	applyTitles(pl, in.Title, in.XLabel, in.YLabel)
	bars, err := plotter.NewBarChart(plotter.Values(in.Values), vg.Points(20))
	if err != nil {
		return nil, output{}, fmt.Errorf("new bar chart: %w", err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	pl.Add(bars)
	pl.NominalX(in.Labels...)
	return p.save(pl, in.Filename, in.Title, "bar")
}

// ── common ───────────────────────────────────────────────────

type output struct {
	Path     string `json:"path"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	BytesB64 string `json:"bytes_b64"`
}

func (p *plugin) save(pl *plot.Plot, filename, title, fallback string) (*mcp.CallToolResult, output, error) {
	if err := os.MkdirAll(p.cfg.OutputDir, 0755); err != nil {
		return nil, output{}, fmt.Errorf("mkdir %s: %w", p.cfg.OutputDir, err)
	}
	name := strings.TrimSpace(filename)
	if name == "" {
		base := slug(title)
		if base == "" {
			base = fallback
		}
		name = fmt.Sprintf("%s_%s.png", base, time.Now().UTC().Format("20060102T150405Z"))
	}
	if !strings.HasSuffix(strings.ToLower(name), ".png") {
		name += ".png"
	}
	path := filepath.Join(p.cfg.OutputDir, name)
	w := vg.Points(float64(p.cfg.Width)) / vg.Inch * vg.Length(p.cfg.DPI) / vg.Length(p.cfg.DPI)
	_ = w // (compute kept for future DPI scaling; gonum sizes in vg.Length)
	if err := pl.Save(vg.Length(p.cfg.Width), vg.Length(p.cfg.Height), path); err != nil {
		return nil, output{}, fmt.Errorf("save plot: %w", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, output{}, fmt.Errorf("read png: %w", err)
	}
	if p.cfg.AutoOpen != nil && *p.cfg.AutoOpen {
		_ = openInViewer(path)
	}
	o := output{
		Path:     path,
		Width:    p.cfg.Width,
		Height:   p.cfg.Height,
		BytesB64: base64.StdEncoding.EncodeToString(raw),
	}
	b, _ := json.Marshal(map[string]any{"path": o.Path, "width": o.Width, "height": o.Height})
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(b)},
			&mcp.ImageContent{Data: raw, MIMEType: "image/png"},
		},
	}, o, nil
}

func applyTitles(pl *plot.Plot, title, xl, yl string) {
	if title != "" {
		pl.Title.Text = title
	}
	if xl != "" {
		pl.X.Label.Text = xl
	}
	if yl != "" {
		pl.Y.Label.Text = yl
	}
}

func toXYs(s Series) plotter.XYs {
	n := len(s.Y)
	xys := make(plotter.XYs, n)
	useX := len(s.X) == n
	for i := 0; i < n; i++ {
		if useX {
			xys[i].X = s.X[i]
		} else {
			xys[i].X = float64(i)
		}
		xys[i].Y = s.Y[i]
	}
	return xys
}

// openInViewer launches the OS default image viewer for path. Best-effort:
// on headless systems or missing helpers it returns an error which the caller
// is expected to ignore (the PNG is still on disk).
func openInViewer(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", abs)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", abs)
	default:
		cmd = exec.Command("xdg-open", abs)
	}
	return cmd.Start()
}

var slugRE = regexp.MustCompile(`[^a-z0-9]+`)

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRE.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 60 {
		s = s[:60]
	}
	return s
}

// ── csv-backed plot tools ────────────────────────────────────

type csvLineInput struct {
	Path      string   `json:"path" jsonschema:"absolute or working-dir-relative path to the CSV file"`
	Title     string   `json:"title,omitempty"`
	XLabel    string   `json:"x_label,omitempty"`
	YLabel    string   `json:"y_label,omitempty"`
	Filename  string   `json:"filename,omitempty"`
	SymbolCol string   `json:"symbol_col,omitempty" jsonschema:"name of the symbol column to filter on (default: sym)"`
	Symbol    string   `json:"symbol,omitempty" jsonschema:"row filter on symbol_col; empty plots every row"`
	XCol      string   `json:"x_col,omitempty" jsonschema:"optional numeric column for the x axis; defaults to row index"`
	YCols     []string `json:"y_cols" jsonschema:"one or more numeric columns to plot; each becomes its own series"`
}

func (p *plugin) handleCSVLine(_ context.Context, _ *mcp.CallToolRequest, in csvLineInput) (*mcp.CallToolResult, output, error) {
	cols := append([]string{}, in.YCols...)
	if in.XCol != "" {
		cols = append(cols, in.XCol)
	}
	data, err := readCSVColumns(in.Path, in.SymbolCol, in.Symbol, cols)
	if err != nil {
		return nil, output{}, err
	}
	series := make([]Series, 0, len(in.YCols))
	for _, name := range in.YCols {
		ys, ok := data[strings.ToLower(name)]
		if !ok {
			return nil, output{}, fmt.Errorf("column %q not found", name)
		}
		s := Series{Name: name, Y: ys}
		if in.XCol != "" {
			s.X = data[strings.ToLower(in.XCol)]
		}
		series = append(series, s)
	}
	pl := plot.New()
	applyTitles(pl, in.Title, in.XLabel, in.YLabel)
	args := make([]any, 0, 2*len(series))
	for _, s := range series {
		args = append(args, s.Name, toXYs(s))
	}
	if err := plotutil.AddLines(pl, args...); err != nil {
		return nil, output{}, fmt.Errorf("add lines: %w", err)
	}
	return p.save(pl, in.Filename, csvTitle(in.Title, in.Path, in.Symbol), "csv_line")
}

func (p *plugin) handleCSVScatter(_ context.Context, _ *mcp.CallToolRequest, in csvLineInput) (*mcp.CallToolResult, output, error) {
	cols := append([]string{}, in.YCols...)
	if in.XCol != "" {
		cols = append(cols, in.XCol)
	}
	data, err := readCSVColumns(in.Path, in.SymbolCol, in.Symbol, cols)
	if err != nil {
		return nil, output{}, err
	}
	series := make([]Series, 0, len(in.YCols))
	for _, name := range in.YCols {
		ys, ok := data[strings.ToLower(name)]
		if !ok {
			return nil, output{}, fmt.Errorf("column %q not found", name)
		}
		s := Series{Name: name, Y: ys}
		if in.XCol != "" {
			s.X = data[strings.ToLower(in.XCol)]
		}
		series = append(series, s)
	}
	pl := plot.New()
	applyTitles(pl, in.Title, in.XLabel, in.YLabel)
	args := make([]any, 0, 2*len(series))
	for _, s := range series {
		args = append(args, s.Name, toXYs(s))
	}
	if err := plotutil.AddScatters(pl, args...); err != nil {
		return nil, output{}, fmt.Errorf("add scatters: %w", err)
	}
	return p.save(pl, in.Filename, csvTitle(in.Title, in.Path, in.Symbol), "csv_scatter")
}

type csvHistInput struct {
	Path      string `json:"path"`
	Title     string `json:"title,omitempty"`
	XLabel    string `json:"x_label,omitempty"`
	YLabel    string `json:"y_label,omitempty"`
	Filename  string `json:"filename,omitempty"`
	SymbolCol string `json:"symbol_col,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Column    string `json:"column" jsonschema:"numeric column to bin"`
	Bins      int    `json:"bins,omitempty"`
}

func (p *plugin) handleCSVHistogram(_ context.Context, _ *mcp.CallToolRequest, in csvHistInput) (*mcp.CallToolResult, output, error) {
	data, err := readCSVColumns(in.Path, in.SymbolCol, in.Symbol, []string{in.Column})
	if err != nil {
		return nil, output{}, err
	}
	values, ok := data[strings.ToLower(in.Column)]
	if !ok || len(values) == 0 {
		return nil, output{}, fmt.Errorf("column %q has no values", in.Column)
	}
	pl := plot.New()
	applyTitles(pl, in.Title, in.XLabel, in.YLabel)
	if in.YLabel == "" {
		pl.Y.Label.Text = "frequency"
	}
	bins := in.Bins
	if bins <= 0 {
		bins = 30
	}
	h, err := plotter.NewHist(plotter.Values(values), bins)
	if err != nil {
		return nil, output{}, fmt.Errorf("new histogram: %w", err)
	}
	h.FillColor = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	pl.Add(h)
	return p.save(pl, in.Filename, csvTitle(in.Title, in.Path, in.Symbol), "csv_histogram")
}

// readCSVColumns opens path, parses the header, and returns each requested
// column as a []float64. If symbolCol+symbol are both non-empty, only matching
// rows are kept (case-insensitive). Column lookup is case-insensitive; map
// keys are returned lowercased.
func readCSVColumns(path, symbolCol, symbol string, cols []string) (map[string][]float64, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("path is required")
	}
	if len(cols) == 0 {
		return nil, fmt.Errorf("at least one column is required")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	r := stdcsv.NewReader(f)
	r.FieldsPerRecord = -1
	hdr, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	idx := map[string]int{}
	for i, h := range hdr {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}
	if symbolCol == "" {
		symbolCol = "sym"
	}
	symIdx, hasSym := idx[strings.ToLower(symbolCol)]
	colIdx := make([]int, len(cols))
	for i, c := range cols {
		j, ok := idx[strings.ToLower(c)]
		if !ok {
			return nil, fmt.Errorf("column %q not in header %v", c, hdr)
		}
		colIdx[i] = j
	}
	target := strings.ToUpper(strings.TrimSpace(symbol))
	out := make(map[string][]float64, len(cols))
	for _, c := range cols {
		out[strings.ToLower(c)] = nil
	}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		if target != "" {
			if !hasSym {
				return nil, fmt.Errorf("symbol filter set but column %q not in header", symbolCol)
			}
			if strings.ToUpper(row[symIdx]) != target {
				continue
			}
		}
		for i, c := range cols {
			v, err := strconv.ParseFloat(row[colIdx[i]], 64)
			if err != nil {
				continue
			}
			key := strings.ToLower(c)
			out[key] = append(out[key], v)
		}
	}
	return out, nil
}

func csvTitle(title, path, sym string) string {
	if title != "" {
		return title
	}
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if sym != "" {
		return base + "-" + sym
	}
	return base
}

func init() { tools.Register(&plugin{}) }
