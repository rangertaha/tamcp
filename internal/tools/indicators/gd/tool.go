// Package gd registers the Generalized DEMA (Pandas TA).
package gd

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "gd",
		Description: "Generalized DEMA: (1+v) * EMA(real, p) - v * EMA(EMA(real, p), p). v=0 → SMA-like; v=1 → DEMA.",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 10},
			{Name: "v", Type: "float", Default: 0.7, Desc: "DEMA volume factor in [0,1]"},
		},
		Run: func(args map[string]any) (any, string, error) {
			val, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 10)
			if p <= 0 {
				p = 10
			}
			v := talib.ArgFloat(args, "v", 0.7)
			out := talib.GDFn(val, p, v)
			return talib.One(out), talib.Tersum("gd", out), nil
		},
	})
}
