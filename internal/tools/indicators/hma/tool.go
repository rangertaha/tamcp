// Package hma registers the Hull Moving Average indicator (Pandas TA).
package hma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "hma",
		Description: "Hull Moving Average (Alan Hull): WMA(2*WMA(p/2) - WMA(p), √p).",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("hma", 20, talib.HMAFn),
	})
}
