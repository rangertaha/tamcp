// Package bbsqueeze registers a classic Bollinger Squeeze boolean.
package bbsqueeze

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "bb_squeeze",
		Description: "Boolean (1/0): 1 when the current Bollinger Bandwidth is at the lowest in the last `lookback` bars.",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 20},
			{Name: "dev", Type: "float", Default: 2.0},
			{Name: "lookback", Type: "int", Default: 120, Desc: "rolling window for the BBW minimum"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			lb := talib.ArgInt(args, "lookback", 120)
			if lb <= 0 {
				lb = 120
			}
			d := talib.ArgFloat(args, "dev", 2.0)
			out := talib.BBSQUEEZEFn(v, p, d, lb)
			return talib.One(out), talib.Tersum("bb_squeeze", out), nil
		},
	})
}
