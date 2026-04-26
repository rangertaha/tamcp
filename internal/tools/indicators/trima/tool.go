// Package trima registers the trima indicator with the talib dispatcher.
package trima

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "trima",
		Description: "Triangular Moving Average",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("trima", 30, talib.TRIMAFn),
	})
}
