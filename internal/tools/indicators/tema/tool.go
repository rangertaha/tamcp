// Package tema registers the tema indicator with the talib dispatcher.
package tema

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tema",
		Description: "Triple Exponential Moving Average",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("tema", 30, talib.TEMAFn),
	})
}
