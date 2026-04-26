// Package adsmooth registers an EMA-smoothed AD line (utility).
package adsmooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ad_smoothed",
		Description: "EMA-smoothed Accumulation/Distribution line.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "period", Type: "int", Default: 14, Desc: "EMA period"},
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
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			out := talib.ADSMOOTHFn(h, l, c, v, p)
			return talib.One(out), talib.Tersum("ad_smoothed", out), nil
		},
	})
}
