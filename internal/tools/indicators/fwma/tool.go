// Package fwma registers the Fibonacci Weighted Moving Average (Pandas TA).
package fwma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "fwma",
		Description: "Fibonacci Weighted MA: weights are F_2..F_(p+1), newest weight largest.",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("fwma", 10, talib.FWMAFn),
	})
}
