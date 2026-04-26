// Package squeezepro registers TTM Squeeze Pro (Pandas TA).
package squeezepro

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "squeeze_pro",
		Description: "TTM Squeeze Pro. Tracks Bollinger-inside-Keltner compression at 3 KC widths plus the LINEARREG momentum series.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "bb_length", Type: "int", Default: 20},
			{Name: "bb_mult", Type: "float", Default: 2.0},
			{Name: "kc_length", Type: "int", Default: 20},
			{Name: "kc_mult_low", Type: "float", Default: 2.0},
			{Name: "kc_mult_mid", Type: "float", Default: 1.5},
			{Name: "kc_mult_high", Type: "float", Default: 1.0},
			{Name: "mom_length", Type: "int", Default: 12},
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
			sLow, sMid, sHigh, mom := talib.SQUEEZEPROFn(h, l, c,
				talib.ArgInt(args, "bb_length", 20),
				talib.ArgFloat(args, "bb_mult", 2.0),
				talib.ArgInt(args, "kc_length", 20),
				talib.ArgFloat(args, "kc_mult_low", 2.0),
				talib.ArgFloat(args, "kc_mult_mid", 1.5),
				talib.ArgFloat(args, "kc_mult_high", 1.0),
				talib.ArgInt(args, "mom_length", 12),
			)
			out := map[string]any{
				"squeeze_low":  sLow,
				"squeeze_mid":  sMid,
				"squeeze_high": sHigh,
				"momentum":     mom,
			}
			return out, talib.Tersum("squeeze_pro", mom), nil
		},
	})
}
