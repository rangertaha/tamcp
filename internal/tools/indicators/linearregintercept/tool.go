// Package linearregintercept registers the linearreg_intercept indicator with the talib dispatcher.
package linearregintercept

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "linearreg_intercept",
		Description: "Linear Regression Intercept",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("linearreg_intercept", 14, talib.LINEARREGINTERCEPTFn),
	})
}
