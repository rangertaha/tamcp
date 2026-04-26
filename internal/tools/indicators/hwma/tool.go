// Package hwma registers the Holt-Winters Moving Average (Pandas TA).
package hwma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "hwma",
		Description: "Holt-Winters MA: triple-exponential smoothing with level, trend, and acceleration constants.",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "na", Type: "float", Default: 0.2, Desc: "level smoothing"},
			{Name: "nb", Type: "float", Default: 0.1, Desc: "trend smoothing"},
			{Name: "nc", Type: "float", Default: 0.1, Desc: "acceleration smoothing"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.HWMAFn(v,
				talib.ArgFloat(args, "na", 0.2),
				talib.ArgFloat(args, "nb", 0.1),
				talib.ArgFloat(args, "nc", 0.1),
			)
			return talib.One(out), talib.Tersum("hwma", out), nil
		},
	})
}
