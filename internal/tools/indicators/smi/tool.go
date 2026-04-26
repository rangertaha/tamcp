// Package smi registers the Stochastic Momentum Index (Pandas TA).
package smi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "smi",
		Description: "Stochastic Momentum Index. Returns smi and an EMA-smoothed signal.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "length", Type: "int", Default: 10, Desc: "stochastic look-back"},
			{Name: "smooth_period", Type: "int", Default: 3, Desc: "double-EMA smoothing"},
			{Name: "signal_period", Type: "int", Default: 3, Desc: "EMA period for signal"},
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
			s, sig := talib.SMIFn(h, l, c,
				talib.ArgInt(args, "length", 10),
				talib.ArgInt(args, "smooth_period", 3),
				talib.ArgInt(args, "signal_period", 3),
			)
			return talib.Two(s, sig, [2]string{"smi", "signal"}), talib.Tersum("smi", s), nil
		},
	})
}
