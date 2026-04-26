// Package pctret registers the Percent Return helper (Pandas TA).
package pctret

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "percent_return",
		Description: "Percent return: 100 * (real[i] - real[i-period]) / real[i-period].",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(1),
		Run:         talib.RunRealPeriod("percent_return", 1, talib.PCTRETFn),
	})
}
