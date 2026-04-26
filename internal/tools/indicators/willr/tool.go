// Package willr registers the willr indicator with the talib dispatcher.
package willr

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "willr",
		Description: "Williams' %R",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("willr", 14, talib.WILLRFn),
	})
}
