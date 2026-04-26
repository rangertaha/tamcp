// Package eom registers Ease of Movement (Pandas TA, cinar).
package eom

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "eom",
		Description: "Ease of Movement: SMA( ((h+l)/2 - prev(h+l)/2) / ((volume/divisor)/(h-l)), period ).",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "period", Type: "int", Default: 14, Desc: "SMA period"},
			{Name: "divisor", Type: "float", Default: 100000000.0, Desc: "volume scaling divisor"},
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
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			div := talib.ArgFloat(args, "divisor", 100000000.0)
			if div == 0 {
				div = 100000000.0
			}
			out := talib.EOMFn(h, l, v, p, div)
			return talib.One(out), talib.Tersum("eom", out), nil
		},
	})
}
