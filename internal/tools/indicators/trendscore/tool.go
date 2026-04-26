// Package trendscore registers Trend Score (count of up vs down bars).
package trendscore

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "trend_score",
		Description: "Sum over the last `period` bars of sign(Δclose). Range [-period, +period].",
		Group:       "trend",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("trend_score", 14, talib.TRENDSCOREFn),
	})
}
