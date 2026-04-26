// Package wma registers the wma indicator with the talib dispatcher.
package wma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "wma",
		Description: "Weighted Moving Average",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("wma", 30, talib.WMAFn),
	})
}
