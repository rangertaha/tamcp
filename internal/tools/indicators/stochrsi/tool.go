// Package stochrsi registers the stochrsi indicator with the talib dispatcher.
package stochrsi

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "stochrsi",
		Description: "Stochastic RSI",
		Group:       "momentum",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 14)
			fk := talib.ArgInt(args, "fast_k_period", 5)
			fd := talib.ArgInt(args, "fast_d_period", 3)
			fdma, err := talib.MaTypeFromString(talib.ArgString(args, "fast_d_matype", ""))
			if err != nil {
				return nil, "", err
			}
			k, d := talib.STOCHRSIFn(v, p, fk, fd, fdma)
			return talib.Two(k, d, [2]string{"k", "d"}), talib.Tersum("stochrsi", k), nil
		},
	})
}
