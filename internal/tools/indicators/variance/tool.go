// Package variance registers the var indicator with the talib dispatcher.
package variance

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "var",
		Description: "Variance",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(5),
		Run:         talib.RunRealPeriod("var", 5, talib.VARFn),
	})
}
