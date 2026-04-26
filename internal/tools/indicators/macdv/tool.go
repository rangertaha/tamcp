// Package macdv registers ATR-normalized MACD-V (Pandas TA).
package macdv

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "macd_v",
		Description: "MACD-V: (EMA(close, fast) - EMA(close, slow)) / ATR(slow) * 100. Returns macd_v, signal, histogram.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "fast_period", Type: "int", Default: 12},
			{Name: "slow_period", Type: "int", Default: 26},
			{Name: "signal_period", Type: "int", Default: 9},
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
			m, s, hi := talib.MACDVFn(h, l, c,
				talib.ArgInt(args, "fast_period", 12),
				talib.ArgInt(args, "slow_period", 26),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Three(m, s, hi, [3]string{"macd_v", "signal", "histogram"}), talib.Tersum("macd_v", m), nil
		},
	})
}
