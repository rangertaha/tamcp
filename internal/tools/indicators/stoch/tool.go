// Package stoch registers the stoch indicator with the talib dispatcher.
package stoch

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "stoch",
		Description: "Stochastic Oscillator (slow). Returns k, d.",
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
			sk := talib.ArgInt(args, "slow_k_period", 3)
			skma, err := talib.MaTypeFromString(talib.ArgString(args, "slow_k_matype", ""))
			if err != nil {
				return nil, "", err
			}
			sd := talib.ArgInt(args, "slow_d_period", 3)
			sdma, err := talib.MaTypeFromString(talib.ArgString(args, "slow_d_matype", ""))
			if err != nil {
				return nil, "", err
			}
			k, d := talib.STOCHFn(h, l, c, fk, sk, skma, sd, sdma)
			return talib.Two(k, d, [2]string{"k", "d"}), talib.Tersum("stoch", k), nil
		},
	})
}
