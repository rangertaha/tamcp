// Package tdi registers Trader's Dynamic Index (Pandas TA).
package tdi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tdi",
		Description: "Trader's Dynamic Index. Returns rsi (raw), mab, mbl, upper, middle, lower bands of RSI.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "rsi_period", Type: "int", Default: 13},
			{Name: "mab_period", Type: "int", Default: 2, Desc: "fast SMA of RSI (market base)"},
			{Name: "mbl_period", Type: "int", Default: 7, Desc: "medium SMA of RSI (signal)"},
			{Name: "band_period", Type: "int", Default: 34, Desc: "Bollinger band period for RSI"},
			{Name: "dev_mult", Type: "float", Default: 1.6185, Desc: "Bollinger stddev multiplier"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			rsi, mab, mbl, u, m, lo := talib.TDIFn(v,
				talib.ArgInt(args, "rsi_period", 13),
				talib.ArgInt(args, "mab_period", 2),
				talib.ArgInt(args, "mbl_period", 7),
				talib.ArgInt(args, "band_period", 34),
				talib.ArgFloat(args, "dev_mult", 1.6185),
			)
			out := map[string]any{
				"rsi": rsi, "mab": mab, "mbl": mbl,
				"upper": u, "middle": m, "lower": lo,
			}
			return out, talib.Tersum("tdi", rsi), nil
		},
	})
}
