// Package pwo registers the Percent Williams Oscillator (Pandas TA).
package pwo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pwo",
		Description: "Percent Williams Oscillator: 100 * (WMA(real, fast) - WMA(real, slow)) / WMA(real, slow).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "fast_period", Type: "int", Default: 13},
			{Name: "slow_period", Type: "int", Default: 34},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.PWOFn(v,
				talib.ArgInt(args, "fast_period", 13),
				talib.ArgInt(args, "slow_period", 34),
			)
			return talib.One(out), talib.Tersum("pwo", out), nil
		},
	})
}
