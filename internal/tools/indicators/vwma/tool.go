// Package vwma registers the Volume Weighted Moving Average (Pandas TA).
package vwma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vwma",
		Description: "Volume Weighted Moving Average: SUM(real*volume, p) / SUM(volume, p).",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series (typically close)"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume series, same length as values"},
			{Name: "period", Type: "int", Default: 20, Desc: "rolling window length"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			vol, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			out := talib.VWMAFn(v, vol, p)
			return talib.One(out), talib.Tersum("vwma", out), nil
		},
	})
}
