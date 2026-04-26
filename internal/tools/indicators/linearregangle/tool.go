// Package linearregangle registers the linearreg_angle indicator with the talib dispatcher.
package linearregangle

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "linearreg_angle",
		Description: "Linear Regression Angle (deg)",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("linearreg_angle", 14, talib.LINEARREGANGLEFn),
	})
}
