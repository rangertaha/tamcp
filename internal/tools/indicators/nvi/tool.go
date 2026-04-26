// Package nvi registers the Negative Volume Index (Pandas TA, cinar).
package nvi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "nvi",
		Description: "Negative Volume Index. Starts at 1000; updates only on bars where volume falls vs the previous bar.",
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
			out := talib.NVIFn(c, v)
			return talib.One(out), talib.Tersum("nvi", out), nil
		},
	})
}
