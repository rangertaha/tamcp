// Package stochdiff registers Slow K - Slow D (Stochastic spread).
package stochdiff

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "stoch_diff",
		Description: "Stochastic spread: slow K - slow D (SMA-smoothed) over the configured periods.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "fastk_period", Type: "int", Default: 14},
			{Name: "slowk_period", Type: "int", Default: 3},
			{Name: "slowd_period", Type: "int", Default: 3},
		},
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			out := talib.STOCHDIFFFn(h, l, c,
				talib.ArgInt(args, "fastk_period", 14),
				talib.ArgInt(args, "slowk_period", 3),
				talib.ArgInt(args, "slowd_period", 3),
			)
			return talib.One(out), talib.Tersum("stoch_diff", out), nil
		},
	})
}
