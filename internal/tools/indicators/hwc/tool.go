// Package hwc registers the Holt-Winters Channel (Pandas TA).
package hwc

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "hwc",
		Description: "Holt-Winters Channel: HWMA mid ± scalar * STDDEV(close - HWMA, channel_period). Returns upper, middle, lower.",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "na", Type: "float", Default: 0.2},
			{Name: "nb", Type: "float", Default: 0.1},
			{Name: "nc", Type: "float", Default: 0.1},
			{Name: "channel_period", Type: "int", Default: 20},
			{Name: "scalar", Type: "float", Default: 1.0},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			u, m, lo := talib.HWCFn(v,
				talib.ArgFloat(args, "na", 0.2),
				talib.ArgFloat(args, "nb", 0.1),
				talib.ArgFloat(args, "nc", 0.1),
				talib.ArgInt(args, "channel_period", 20),
				talib.ArgFloat(args, "scalar", 1.0),
			)
			return talib.Three(u, m, lo, [3]string{"upper", "middle", "lower"}), talib.Tersum("hwc", m), nil
		},
	})
}
