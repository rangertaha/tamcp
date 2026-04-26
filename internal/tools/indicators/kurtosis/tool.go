// Package kurtosis registers the rolling kurtosis indicator (Pandas TA).
package kurtosis

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "kurtosis",
		Description: "Rolling sample excess kurtosis over `period` bars: m4/m2² - 3.",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("kurtosis", 30, talib.KURTOSISFn),
	})
}
