// Package tsf registers the tsf indicator with the talib dispatcher.
package tsf

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tsf",
		Description: "Time Series Forecast",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("tsf", 14, talib.TSFFn),
	})
}
