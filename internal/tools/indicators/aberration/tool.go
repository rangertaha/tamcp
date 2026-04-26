// Package aberration registers the Aberration ATR-banded SMA (Pandas TA).
package aberration

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "aberration",
		Description: "Aberration: ATR bands around SMA(HLC3). Returns upper, middle, lower.",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "length", Type: "int", Default: 5, Desc: "SMA period for HLC3 mid line"},
			{Name: "atr_length", Type: "int", Default: 15, Desc: "ATR period for the bands"},
			{Name: "atr_mult", Type: "float", Default: 1.0, Desc: "ATR band multiplier"},
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
			u, m, lo := talib.ABERRATIONFn(h, l, c,
				talib.ArgInt(args, "length", 5),
				talib.ArgInt(args, "atr_length", 15),
				talib.ArgFloat(args, "atr_mult", 1.0),
			)
			return talib.Three(u, m, lo, [3]string{"upper", "middle", "lower"}), talib.Tersum("aberration", m), nil
		},
	})
}
