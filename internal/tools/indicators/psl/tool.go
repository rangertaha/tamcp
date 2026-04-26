// Package psl registers the Psychological Line indicator (Pandas TA).
package psl

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "psl",
		Description: "Psychological Line: 100 * (count of up bars over the last p bars) / p.",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(12),
		Run:         talib.RunRealPeriod("psl", 12, talib.PSLFn),
	})
}
