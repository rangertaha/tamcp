// Package plusdm registers the plus_dm indicator with the talib dispatcher.
package plusdm

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "plus_dm",
		Description: "Plus Directional Movement",
		Group:       "momentum",
		Params:      talib.ParamsHLPeriod(14),
		Run:         talib.RunHLPeriod("plus_dm", 14, talib.PLUSDMFn),
	})
}
