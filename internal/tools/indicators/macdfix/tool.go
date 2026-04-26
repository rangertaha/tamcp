// Package macdfix registers the macdfix indicator with the talib dispatcher.
package macdfix

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "macdfix",
		Description: "MACDFIX (fast=12, slow=26 fixed).",
		Group:       "momentum",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			gp := talib.ArgInt(args, "signal_period", 9)
			m, sg, h := talib.MACDFIXFn(v, gp)
			return talib.Three(m, sg, h, [3]string{"macd", "signal", "histogram"}), talib.Tersum("macdfix", m), nil
		},
	})
}
