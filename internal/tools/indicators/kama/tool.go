// Package kama registers the kama indicator with the talib dispatcher.
package kama

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "kama",
		Description: "Kaufman Adaptive Moving Average",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("kama", 30, talib.KAMAFn),
	})
}
