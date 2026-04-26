// Package ttmtrend registers the TTM Trend indicator (Pandas TA).
package ttmtrend

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ttm_trend",
		Description: "TTM Trend bar: +1 if close > SMA((H+L)/2, period), else -1.",
		Group:       "trend",
		Params:      talib.ParamsHLCPeriod(6),
		Run:         talib.RunHLCPeriod("ttm_trend", 6, talib.TTMTRENDFn),
	})
}
