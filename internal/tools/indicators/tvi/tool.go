// Package tvi registers the Trade Volume Index (cinar).
package tvi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tvi",
		Description: "Trade Volume Index: cumulative sign(Δclose vs ±min_tick) * volume.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "min_tick", Type: "float", Default: 0.5, Desc: "minimum price move that counts as a direction change"},
		},
		Run: func(args map[string]any) (any, string, error) {
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			t := talib.ArgFloat(args, "min_tick", 0.5)
			out := talib.TVIFn(c, v, t)
			return talib.One(out), talib.Tersum("tvi", out), nil
		},
	})
}
