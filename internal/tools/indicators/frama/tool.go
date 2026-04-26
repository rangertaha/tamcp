// Package frama registers Ehlers' Fractal Adaptive Moving Average (Pandas TA).
package frama

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "frama",
		Description: "Fractal Adaptive Moving Average (Ehlers): EMA whose smoothing constant adapts to the fractal dimension over an even `period`.",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 16, Desc: "even period; rounded down to the nearest even"},
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
			p := talib.ArgInt(args, "period", 16)
			if p < 2 {
				p = 16
			}
			if p%2 != 0 {
				p--
			}
			out := talib.FRAMAFn(h, l, c, p)
			return talib.One(out), talib.Tersum("frama", out), nil
		},
	})
}
