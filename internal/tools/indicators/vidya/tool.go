// Package vidya registers Chande's Variable Index Dynamic Average (Pandas TA, cinar).
package vidya

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vidya",
		Description: "Variable Index Dynamic Average: EMA whose smoothing is scaled by |CMO|/100 over `cmo_period`.",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 14, Desc: "EMA-equivalent period"},
			{Name: "cmo_period", Type: "int", Default: 9, Desc: "CMO look-back used as the volatility factor"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			cp := talib.ArgInt(args, "cmo_period", 9)
			if cp <= 0 {
				cp = 9
			}
			out := talib.VIDYAFn(v, p, cp)
			return talib.One(out), talib.Tersum("vidya", out), nil
		},
	})
}
