// Package wmadiff registers WMA(fast) - WMA(slow) (utility).
package wmadiff

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "wma_diff",
		Description: "WMA(real, fast) - WMA(real, slow). WMA-based MACD-line analogue.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "fast_period", Type: "int", Default: 12},
			{Name: "slow_period", Type: "int", Default: 26},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.WMADIFFFn(v,
				talib.ArgInt(args, "fast_period", 12),
				talib.ArgInt(args, "slow_period", 26),
			)
			return talib.One(out), talib.Tersum("wma_diff", out), nil
		},
	})
}
