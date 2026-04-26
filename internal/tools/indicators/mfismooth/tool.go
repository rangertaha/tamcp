// Package mfismooth registers an EMA-smoothed Money Flow Index.
package mfismooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mfi_smoothed",
		Description: "EMA-smoothed Money Flow Index.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "mfi_period", Type: "int", Default: 14},
			{Name: "smooth_period", Type: "int", Default: 5},
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
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			mp := talib.ArgInt(args, "mfi_period", 14)
			sp := talib.ArgInt(args, "smooth_period", 5)
			out := talib.MFISMOOTHFn(h, l, c, v, mp, sp)
			return talib.One(out), talib.Tersum("mfi_smoothed", out), nil
		},
	})
}
