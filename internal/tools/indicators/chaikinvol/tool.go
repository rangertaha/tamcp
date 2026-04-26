// Package chaikinvol registers the Chaikin Volatility indicator (cinar).
package chaikinvol

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "chaikin_vol",
		Description: "Chaikin Volatility: 100 * ROC(EMA(H-L, ema_period), roc_period).",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "ema_period", Type: "int", Default: 10, Desc: "EMA window over the H-L series"},
			{Name: "roc_period", Type: "int", Default: 10, Desc: "ROC look-back over the EMA series"},
		},
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			out := talib.CHAIKINVOLFn(h, l,
				talib.ArgInt(args, "ema_period", 10),
				talib.ArgInt(args, "roc_period", 10),
			)
			return talib.One(out), talib.Tersum("chaikin_vol", out), nil
		},
	})
}
