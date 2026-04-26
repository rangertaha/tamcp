// Package sumwindow registers the sum_window indicator with the talib dispatcher.
package sumwindow

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sum_window",
		Description: "Rolling sum (TA-Lib SUM)",
		Group:       "operator",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("sum_window", 30, talib.SUMWINDOWFn),
	})
}
