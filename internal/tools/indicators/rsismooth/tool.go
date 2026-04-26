// Package rsismooth registers an EMA-smoothed RSI (utility, Pandas TA-style).
package rsismooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rsi_smoothed",
		Description: "EMA-smoothed RSI: EMA(RSI(close, rsi_period), smooth_period).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "rsi_period", Type: "int", Default: 14},
			{Name: "smooth_period", Type: "int", Default: 5},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			rp := talib.ArgInt(args, "rsi_period", 14)
			if rp <= 0 {
				rp = 14
			}
			sp := talib.ArgInt(args, "smooth_period", 5)
			if sp <= 0 {
				sp = 5
			}
			out := talib.RSISMOOTHFn(v, rp, sp)
			return talib.One(out), talib.Tersum("rsi_smoothed", out), nil
		},
	})
}
