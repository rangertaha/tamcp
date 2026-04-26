// Package alma registers the Arnaud Legoux Moving Average (Pandas TA, cinar).
package alma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "alma",
		Description: "Arnaud Legoux Moving Average. Gaussian-weighted MA with `offset` (0..1) and `sigma` shape parameters.",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "window", Type: "int", Default: 9, Desc: "window size"},
			{Name: "sigma", Type: "float", Default: 6.0, Desc: "Gaussian width (smaller = smoother)"},
			{Name: "offset", Type: "float", Default: 0.85, Desc: "centre of weighting in [0,1]"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			w := talib.ArgInt(args, "window", 9)
			if w <= 0 {
				w = 9
			}
			sigma := talib.ArgFloat(args, "sigma", 6.0)
			if sigma <= 0 {
				sigma = 6.0
			}
			off := talib.ArgFloat(args, "offset", 0.85)
			out := talib.ALMAFn(v, w, sigma, off)
			return talib.One(out), talib.Tersum("alma", out), nil
		},
	})
}
