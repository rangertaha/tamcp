// Package pvi registers the Positive Volume Index (Pandas TA, cinar).
package pvi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pvi",
		Description: "Positive Volume Index. Starts at 1000; updates only on bars where volume rises vs the previous bar.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
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
			out := talib.PVIFn(c, v)
			return talib.One(out), talib.Tersum("pvi", out), nil
		},
	})
}
