// Package minindex registers the minindex indicator with the talib dispatcher.
package minindex

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "minindex",
		Description: "Index of rolling minimum",
		Group:       "operator",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("minindex", 30, talib.MININDEXFn),
	})
}
