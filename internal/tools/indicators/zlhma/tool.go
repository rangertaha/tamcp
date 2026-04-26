// Package zlhma registers the Zero-Lag Hull MA (Pandas TA).
package zlhma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "zlhma",
		Description: "Zero-Lag Hull MA: HMA computed on a de-lagged source (src[i] = 2*real[i] - real[i-(p-1)/2]).",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("zlhma", 20, talib.ZLHMAFn),
	})
}
