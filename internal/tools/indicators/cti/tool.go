// Package cti registers the Correlation Trend Indicator (Pandas TA).
package cti

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cti",
		Description: "Correlation Trend Indicator: rolling Pearson correlation of price vs a linear time index over p bars.",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(12),
		Run:         talib.RunRealPeriod("cti", 12, talib.CTIFn),
	})
}
