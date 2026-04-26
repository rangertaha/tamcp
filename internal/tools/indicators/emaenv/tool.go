// Package emaenv registers an EMA-based Envelope (variant of Pandas TA's envelope).
package emaenv

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ema_envelope",
		Description: "EMA-based Envelope: EMA(close, p) * (1 ± pct/100). Returns upper, middle, lower.",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 20},
			{Name: "percent", Type: "float", Default: 2.5},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			pct := talib.ArgFloat(args, "percent", 2.5)
			u, m, lo := talib.EMAENVFn(v, p, pct)
			return talib.Three(u, m, lo, [3]string{"upper", "middle", "lower"}), talib.Tersum("ema_envelope", m), nil
		},
	})
}
