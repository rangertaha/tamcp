// Package linearreg registers the linearreg indicator with the talib dispatcher.
package linearreg

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "linearreg",
		Description: "Linear Regression",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("linearreg", 14, talib.LINEARREGFn),
	})
}
