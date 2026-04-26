// Package macd registers the macd indicator with the talib dispatcher.
package macd

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "macd",
		Description: "MACD. Returns macd, signal, histogram.",
		Group:       "momentum",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			fp := talib.ArgInt(args, "fast_period", 12)
			sp := talib.ArgInt(args, "slow_period", 26)
			gp := talib.ArgInt(args, "signal_period", 9)
			m, sg, h := talib.MACDFn(v, fp, sp, gp)
			return talib.Three(m, sg, h, [3]string{"macd", "signal", "histogram"}), talib.Tersum("macd", m), nil
		},
	})
}
