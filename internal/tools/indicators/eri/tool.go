// Package eri registers Elder Ray Index (bull / bear power) (Pandas TA).
package eri

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "eri",
		Description: "Elder Ray Index. Returns bull = high - EMA(close,p) and bear = low - EMA(close,p).",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(13),
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 13)
			if p <= 0 {
				p = 13
			}
			bull, bear := talib.ERIFn(h, l, c, p)
			return talib.Two(bull, bear, [2]string{"bull", "bear"}), talib.Tersum("eri", bull), nil
		},
	})
}
