// Package macdzl registers the Zero-Lag MACD (Pandas TA-style).
package macdzl

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "macd_zero_lag",
		Description: "Zero-Lag MACD: ZLEMA(fast) - ZLEMA(slow), with ZLEMA signal and histogram.",
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
			m, s, h := talib.MACDZLFn(v,
				talib.ArgInt(args, "fast_period", 12),
				talib.ArgInt(args, "slow_period", 26),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Three(m, s, h, [3]string{"macd", "signal", "histogram"}), talib.Tersum("macd_zero_lag", m), nil
		},
	})
}
