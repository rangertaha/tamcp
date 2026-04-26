// Package mcginley registers the McGinley Dynamic indicator (Pandas TA).
package mcginley

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mcginley",
		Description: "McGinley Dynamic: MD[i] = MD[i-1] + (c - MD[i-1]) / (k * p * (c/MD[i-1])^4).",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 10, Desc: "smoothing period"},
			{Name: "k", Type: "float", Default: 0.6, Desc: "responsiveness constant"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 10)
			if p <= 0 {
				p = 10
			}
			k := talib.ArgFloat(args, "k", 0.6)
			if k <= 0 {
				k = 0.6
			}
			out := talib.MCGINLEYFn(v, p, k)
			return talib.One(out), talib.Tersum("mcginley", out), nil
		},
	})
}
