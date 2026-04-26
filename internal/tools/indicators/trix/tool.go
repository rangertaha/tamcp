// Package trix registers the trix indicator with the talib dispatcher.
package trix

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "trix",
		Description: "1-day ROC of triple-smoothed EMA",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("trix", 30, talib.TRIXFn),
	})
}
