// Package dema registers the dema indicator with the talib dispatcher.
package dema

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "dema",
		Description: "Double Exponential Moving Average",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("dema", 30, talib.DEMAFn),
	})
}
