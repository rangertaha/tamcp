// Package increasing registers the Increasing boolean indicator (Pandas TA).
package increasing

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "increasing",
		Description: "Boolean (1/0) — close[i] > close[i-length].",
		Group:       "trend",
		Params:      talib.ParamsRealPeriod(1),
		Run:         talib.RunRealPeriod("increasing", 1, talib.INCREASINGFn),
	})
}
