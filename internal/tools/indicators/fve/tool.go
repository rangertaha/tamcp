// Package fve registers the Finite Volume Element (Markos Katsanos / Pandas TA).
package fve

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "fve",
		Description: "Finite Volume Element: directional volume index using a STDDEV-based intraday cutoff.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "period", Type: "int", Default: 22, Desc: "rolling window for cutoff and sums"},
			{Name: "factor", Type: "float", Default: 0.3, Desc: "intra-day cutoff scaling factor"},
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
			p := talib.ArgInt(args, "period", 22)
			if p <= 0 {
				p = 22
			}
			f := talib.ArgFloat(args, "factor", 0.3)
			out := talib.FVEFn(h, l, c, v, p, f)
			return talib.One(out), talib.Tersum("fve", out), nil
		},
	})
}
