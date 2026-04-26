// Package brar registers the BRAR indicator (Pandas TA).
package brar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "brar",
		Description: "BRAR. Returns AR (open vs range) and BR (prior close vs range) over `period` bars.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "open", Type: "number[]", Required: true, Desc: "Open prices"},
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 26},
		},
		Run: func(args map[string]any) (any, string, error) {
			o, err := talib.ArgFloats(args, "open")
			if err != nil {
				return nil, "", err
			}
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
			p := talib.ArgInt(args, "period", 26)
			if p <= 0 {
				p = 26
			}
			ar, br := talib.BRARFn(o, h, l, c, p)
			return talib.Two(ar, br, [2]string{"ar", "br"}), talib.Tersum("brar", ar), nil
		},
	})
}
