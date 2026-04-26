// Package aobv registers Accumulation On-Balance Volume (cinar).
package aobv

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "aobv",
		Description: "Accumulation OBV: fast/slow EMAs over OBV for crossover signals.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "fast_period", Type: "int", Default: 4},
			{Name: "slow_period", Type: "int", Default: 12},
		},
		Run: func(args map[string]any) (any, string, error) {
			c, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			lo, sh := talib.AOBVFn(c, v,
				talib.ArgInt(args, "fast_period", 4),
				talib.ArgInt(args, "slow_period", 12),
			)
			return talib.Two(lo, sh, [2]string{"long", "short"}), talib.Tersum("aobv", lo), nil
		},
	})
}
