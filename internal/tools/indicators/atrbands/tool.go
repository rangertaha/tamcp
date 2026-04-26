// Package atrbands registers ATR Bands (SMA mid ± mult * ATR).
package atrbands

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "atr_bands",
		Description: "ATR Bands: SMA(close, p) ± mult * ATR(h,l,c, p). Returns upper, middle, lower.",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 14},
			{Name: "multiplier", Type: "float", Default: 2.0},
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
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			m := talib.ArgFloat(args, "multiplier", 2.0)
			if m <= 0 {
				m = 2.0
			}
			u, mid, lo := talib.ATRBANDSFn(h, l, c, p, m)
			return talib.Three(u, mid, lo, [3]string{"upper", "middle", "lower"}), talib.Tersum("atr_bands", mid), nil
		},
	})
}
