// Package cksp registers Chande Kroll Stop (Pandas TA).
package cksp

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "chande_kroll",
		Description: "Chande Kroll Stop. Returns long_stop = highest(high - mult*ATR, max_period) and short_stop = lowest(low + mult*ATR, max_period).",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "atr_period", Type: "int", Default: 10},
			{Name: "multiplier", Type: "float", Default: 1.0},
			{Name: "max_period", Type: "int", Default: 9},
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
			lo, sh := talib.CKSPFn(h, l, c,
				talib.ArgInt(args, "atr_period", 10),
				talib.ArgFloat(args, "multiplier", 1.0),
				talib.ArgInt(args, "max_period", 9),
			)
			return talib.Two(lo, sh, [2]string{"long_stop", "short_stop"}), talib.Tersum("chande_kroll", lo), nil
		},
	})
}
