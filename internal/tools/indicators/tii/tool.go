// Package tii registers the Trend Intensity Index (Pandas TA).
package tii

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tii",
		Description: "Trend Intensity Index: 100 * SUM(positive deviations) / (SUM(positive) + SUM(negative)) over (sma_period, sum_period).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "sma_period", Type: "int", Default: 60, Desc: "trend baseline SMA"},
			{Name: "sum_period", Type: "int", Default: 30, Desc: "deviation sum window"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.TIIFn(v,
				talib.ArgInt(args, "sma_period", 60),
				talib.ArgInt(args, "sum_period", 30),
			)
			return talib.One(out), talib.Tersum("tii", out), nil
		},
	})
}
