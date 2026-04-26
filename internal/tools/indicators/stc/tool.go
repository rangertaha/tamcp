// Package stc registers the Schaff Trend Cycle (Pandas TA).
package stc

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "stc",
		Description: "Schaff Trend Cycle: double-stochastic of MACD. Defaults fast=23, slow=50, cycle=10, factor=0.5.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "fast_period", Type: "int", Default: 23},
			{Name: "slow_period", Type: "int", Default: 50},
			{Name: "cycle", Type: "int", Default: 10},
			{Name: "factor", Type: "float", Default: 0.5},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.STCFn(v,
				talib.ArgInt(args, "fast_period", 23),
				talib.ArgInt(args, "slow_period", 50),
				talib.ArgInt(args, "cycle", 10),
				talib.ArgFloat(args, "factor", 0.5),
			)
			return talib.One(out), talib.Tersum("stc", out), nil
		},
	})
}
