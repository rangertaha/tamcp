// Package linearregslope registers the linearreg_slope indicator with the talib dispatcher.
package linearregslope

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "linearreg_slope",
		Description: "Linear Regression Slope",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("linearreg_slope", 14, talib.LINEARREGSLOPEFn),
	})
}
