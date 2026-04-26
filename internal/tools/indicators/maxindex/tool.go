// Package maxindex registers the maxindex indicator with the talib dispatcher.
package maxindex

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "maxindex",
		Description: "Index of rolling maximum",
		Group:       "operator",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("maxindex", 30, talib.MAXINDEXFn),
	})
}
