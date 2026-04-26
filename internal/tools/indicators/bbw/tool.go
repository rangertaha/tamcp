// Package bbw registers Bollinger Bandwidth (Pandas TA).
package bbw

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "bbw",
		Description: "Bollinger Bandwidth (%): 100 * (upper - lower) / middle.",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 20},
			{Name: "nbdevup", Type: "float", Default: 2.0},
			{Name: "nbdevdn", Type: "float", Default: 2.0},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.BBWFn(v,
				talib.ArgInt(args, "period", 20),
				talib.ArgFloat(args, "nbdevup", 2.0),
				talib.ArgFloat(args, "nbdevdn", 2.0),
			)
			return talib.One(out), talib.Tersum("bbw", out), nil
		},
	})
}
