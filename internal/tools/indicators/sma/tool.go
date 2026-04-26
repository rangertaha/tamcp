// Package sma registers the sma indicator with the talib dispatcher.
package sma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sma",
		Description: "Simple Moving Average",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("sma", 30, talib.SMAFn),
	})
}
