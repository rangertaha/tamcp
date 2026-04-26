// Package zscore registers the rolling Z-score indicator (Pandas TA).
package zscore

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "zscore",
		Description: "Rolling Z-score: (real - SMA(real,p)) / STDDEV(real,p).",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("zscore", 30, talib.ZSCOREFn),
	})
}
