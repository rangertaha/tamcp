// Package er registers Kaufman's Efficiency Ratio (Pandas TA).
package er

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "er",
		Description: "Kaufman Efficiency Ratio: |close - close[-p]| / SUM(|Δclose|, p). Range [0,1].",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("er", 10, talib.ERFn),
	})
}
