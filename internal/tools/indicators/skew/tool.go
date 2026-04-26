// Package skew registers the rolling skew indicator (Pandas TA).
package skew

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "skew",
		Description: "Rolling sample skewness over `period` bars: m3 / m2^(3/2).",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("skew", 30, talib.SKEWFn),
	})
}
