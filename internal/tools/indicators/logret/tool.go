// Package logret registers the Log Return helper (Pandas TA).
package logret

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "log_return",
		Description: "Natural log return: ln(real[i] / real[i-period]).",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(1),
		Run:         talib.RunRealPeriod("log_return", 1, talib.LOGRETFn),
	})
}
