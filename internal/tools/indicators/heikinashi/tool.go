// Package heikinashi registers the Heikin-Ashi candle helper (Pandas TA).
package heikinashi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "heikin_ashi",
		Description: "Heikin-Ashi candles. Returns ha_open, ha_high, ha_low, ha_close.",
		Group:       "price",
		Params:      talib.ParamsOHLC(),
		Run: func(args map[string]any) (any, string, error) {
			o, err := talib.ArgFloats(args, "open")
			if err != nil {
				return nil, "", err
			}
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
			ho, hh, hl, hc := talib.HEIKINASHIFn(o, h, l, c)
			out := map[string]any{
				"ha_open":  ho,
				"ha_high":  hh,
				"ha_low":   hl,
				"ha_close": hc,
			}
			return out, talib.Tersum("heikin_ashi", hc), nil
		},
	})
}
