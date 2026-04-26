// Package linregresid registers residuals of a rolling linear regression (utility).
package linregresid

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "linreg_residuals",
		Description: "Residuals from a rolling linear regression: real - LINREG(real, p).",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("linreg_residuals", 20, talib.LINREGRESIDFn),
	})
}
