// Package mad registers Mean Absolute Deviation (Pandas TA).
package mad

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mad",
		Description: "Mean Absolute Deviation around the rolling mean over `period` bars.",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("mad", 30, talib.MADFn),
	})
}
