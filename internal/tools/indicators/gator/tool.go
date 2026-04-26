// Package gator registers the Gator Oscillator (Pandas TA).
package gator

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "gator",
		Description: "Gator Oscillator (Bill Williams): upper = |jaw - teeth|, lower = -|teeth - lips|.",
		Group:       "momentum",
		Params:      talib.ParamsHL(),
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			u, lo := talib.GATORFn(h, l)
			return talib.Two(u, lo, [2]string{"upper", "lower"}), talib.Tersum("gator", u), nil
		},
	})
}
