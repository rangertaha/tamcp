// Package reflex registers Ehlers' Reflex indicator (Pandas TA).
package reflex

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "reflex",
		Description: "Ehlers Reflex: SSF-prefiltered mean-reversion oscillator. Sign-of-momentum companion to TrendFlex.",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("reflex", 20, talib.REFLEXFn),
	})
}
