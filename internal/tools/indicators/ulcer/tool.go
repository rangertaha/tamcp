// Package ulcer registers the Ulcer Index (Pandas TA).
package ulcer

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ulcer",
		Description: "Ulcer Index: sqrt(SUM(drawdown^2, p) / p) where drawdown = 100*(close - max(close,p))/max(close,p).",
		Group:       "volatility",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("ulcer", 14, talib.UIFn),
	})
}
