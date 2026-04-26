// Package smma registers the Smoothed Moving Average / Wilder RMA (Pandas TA, cinar).
package smma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "smma",
		Description: "Smoothed Moving Average (a.k.a. Wilder's RMA): (prev*(p-1)+curr)/p, seeded with the SMA of the first p bars.",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("smma", 14, talib.SMMAFn),
	})
}
