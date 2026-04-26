// Package macdhist registers a MACD-histogram-only indicator (utility for plotting).
package macdhist

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "macd_hist",
		Description: "MACD histogram only (macd - signal). Convenience wrapper for plotting just the histogram.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "fast_period", Type: "int", Default: 12},
			{Name: "slow_period", Type: "int", Default: 26},
			{Name: "signal_period", Type: "int", Default: 9},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.MACDHISTFn(v,
				talib.ArgInt(args, "fast_period", 12),
				talib.ArgInt(args, "slow_period", 26),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.One(out), talib.Tersum("macd_hist", out), nil
		},
	})
}
