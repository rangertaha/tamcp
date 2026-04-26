// Package coppock registers the Coppock Curve (Pandas TA).
package coppock

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "coppock",
		Description: "Coppock Curve: WMA(ROC(c, long) + ROC(c, short), wma_period). Defaults 14/11/10.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "long_period", Type: "int", Default: 14},
			{Name: "short_period", Type: "int", Default: 11},
			{Name: "wma_period", Type: "int", Default: 10},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.COPPOCKFn(v,
				talib.ArgInt(args, "long_period", 14),
				talib.ArgInt(args, "short_period", 11),
				talib.ArgInt(args, "wma_period", 10),
			)
			return talib.One(out), talib.Tersum("coppock", out), nil
		},
	})
}
