// Package ema registers the ema indicator with the talib dispatcher.
package ema

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ema",
		Description: "Exponential Moving Average",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("ema", 30, talib.EMAFn),
	})
}
