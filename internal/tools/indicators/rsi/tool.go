// Package rsi registers the rsi indicator with the talib dispatcher.
package rsi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rsi",
		Description: "Relative Strength Index",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("rsi", 14, talib.RSIFn),
	})
}
