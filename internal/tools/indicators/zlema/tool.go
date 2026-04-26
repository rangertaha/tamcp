// Package zlema registers the Zero-Lag EMA indicator (Pandas TA).
package zlema

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "zlema",
		Description: "Zero-Lag EMA (Ehlers): EMA of price pre-de-lagged by ⌊(p-1)/2⌋ samples.",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("zlema", 20, talib.ZLEMAFn),
	})
}
