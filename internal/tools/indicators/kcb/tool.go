// Package kcb registers Keltner %B (Pandas TA).
package kcb

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "kcb",
		Description: "Keltner %B: (close - lower_kc) / (upper_kc - lower_kc).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 20},
			{Name: "multiplier", Type: "float", Default: 2.0},
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
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			m := talib.ArgFloat(args, "multiplier", 2.0)
			if m <= 0 {
				m = 2.0
			}
			out := talib.KCBFn(h, l, c, p, m)
			return talib.One(out), talib.Tersum("kcb", out), nil
		},
	})
}
