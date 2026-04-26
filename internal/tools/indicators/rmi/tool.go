// Package rmi registers the Relative Momentum Index (Pandas TA).
package rmi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rmi",
		Description: "Relative Momentum Index: RSI variant comparing close to close[t-momentum] with Wilder smoothing.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 14, Desc: "Wilder smoothing period"},
			{Name: "momentum", Type: "int", Default: 4, Desc: "look-back gap for momentum"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			m := talib.ArgInt(args, "momentum", 4)
			if m <= 0 {
				m = 4
			}
			out := talib.RMIFn(v, p, m)
			return talib.One(out), talib.Tersum("rmi", out), nil
		},
	})
}
