// Package roc registers the roc indicator with the talib dispatcher.
package roc

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "roc",
		Description: "Rate of change",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("roc", 10, talib.ROCFn),
	})
}
