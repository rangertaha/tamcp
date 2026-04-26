// Package squeeze registers LazyBear's Squeeze Momentum (Pandas TA).
package squeeze

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "squeeze",
		Description: "LazyBear Squeeze Momentum. Returns squeeze (1 when Bollinger inside Keltner) and momentum (LINEARREG of detrended price).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "bb_length", Type: "int", Default: 20},
			{Name: "bb_mult", Type: "float", Default: 2.0},
			{Name: "kc_length", Type: "int", Default: 20},
			{Name: "kc_mult", Type: "float", Default: 1.5},
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
			sq, mom := talib.SQUEEZEFn(h, l, c,
				talib.ArgInt(args, "bb_length", 20),
				talib.ArgFloat(args, "bb_mult", 2.0),
				talib.ArgInt(args, "kc_length", 20),
				talib.ArgFloat(args, "kc_mult", 1.5),
				talib.ArgInt(args, "mom_length", 12),
			)
			return talib.Two(sq, mom, [2]string{"squeeze", "momentum"}), talib.Tersum("squeeze", mom), nil
		},
	})
}
