// Package kst registers the Know Sure Thing oscillator (Pandas TA).
package kst

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "kst",
		Description: "Know Sure Thing: weighted sum of four ROC SMAs plus an SMA-smoothed signal.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "roc1", Type: "int", Default: 10},
			{Name: "roc2", Type: "int", Default: 15},
			{Name: "roc3", Type: "int", Default: 20},
			{Name: "roc4", Type: "int", Default: 30},
			{Name: "sma1", Type: "int", Default: 10},
			{Name: "sma2", Type: "int", Default: 10},
			{Name: "sma3", Type: "int", Default: 10},
			{Name: "sma4", Type: "int", Default: 15},
			{Name: "signal_period", Type: "int", Default: 9},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			k, sig := talib.KSTFn(v,
				talib.ArgInt(args, "roc1", 10), talib.ArgInt(args, "roc2", 15),
				talib.ArgInt(args, "roc3", 20), talib.ArgInt(args, "roc4", 30),
				talib.ArgInt(args, "sma1", 10), talib.ArgInt(args, "sma2", 10),
				talib.ArgInt(args, "sma3", 10), talib.ArgInt(args, "sma4", 15),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(k, sig, [2]string{"kst", "signal"}), talib.Tersum("kst", k), nil
		},
	})
}
