// Package stochf registers the stochf indicator with the talib dispatcher.
package stochf

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "stochf",
		Description: "Stochastic Fast. Returns k, d.",
		Group:       "momentum",
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
			fk := talib.ArgInt(args, "fast_k_period", 5)
			fd := talib.ArgInt(args, "fast_d_period", 3)
			fdma, err := talib.MaTypeFromString(talib.ArgString(args, "fast_d_matype", ""))
			if err != nil {
				return nil, "", err
			}
			k, d := talib.STOCHFFn(h, l, c, fk, fd, fdma)
			return talib.Two(k, d, [2]string{"k", "d"}), talib.Tersum("stochf", k), nil
		},
	})
}
