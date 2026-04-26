// Package tsi registers the True Strength Index (Pandas TA, cinar).
package tsi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tsi",
		Description: "True Strength Index. Returns tsi and an EMA-smoothed signal.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series (typically close)"},
			{Name: "long_period", Type: "int", Default: 25, Desc: "first EMA smoothing"},
			{Name: "short_period", Type: "int", Default: 13, Desc: "second EMA smoothing"},
			{Name: "signal_period", Type: "int", Default: 13, Desc: "EMA period for signal line"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			r := talib.ArgInt(args, "long_period", 25)
			s := talib.ArgInt(args, "short_period", 13)
			sg := talib.ArgInt(args, "signal_period", 13)
			t, sig := talib.TSIFn(v, r, s, sg)
			return talib.Two(t, sig, [2]string{"tsi", "signal"}), talib.Tersum("tsi", t), nil
		},
	})
}
