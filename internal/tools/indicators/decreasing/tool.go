// Package decreasing registers the Decreasing boolean indicator (Pandas TA).
package decreasing

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "decreasing",
		Description: "Boolean (1/0) — close[i] < close[i-length].",
		Group:       "trend",
		Params:      talib.ParamsRealPeriod(1),
		Run:         talib.RunRealPeriod("decreasing", 1, talib.DECREASINGFn),
	})
}
