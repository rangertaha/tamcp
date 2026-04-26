// Package rocs registers a Smoothed ROC indicator (Pandas TA).
package rocs

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rocs",
		Description: "Smoothed ROC: SMA(ROC(real, period), smooth).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Input series"},
			{Name: "period", Type: "int", Default: 10, Desc: "ROC look-back"},
			{Name: "smooth", Type: "int", Default: 10, Desc: "SMA smoothing on top of ROC"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 10)
			if p <= 0 {
				p = 10
			}
			s := talib.ArgInt(args, "smooth", 10)
			if s <= 0 {
				s = 10
			}
			out := talib.ROCSFn(v, p, s)
			return talib.One(out), talib.Tersum("rocs", out), nil
		},
	})
}
