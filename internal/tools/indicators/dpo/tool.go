// Package dpo registers the Detrended Price Oscillator (Pandas TA).
package dpo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "dpo",
		Description: "Detrended Price Oscillator: close[i - (p/2+1)] - SMA(close, p)[i].",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("dpo", 20, talib.DPOFn),
	})
}
