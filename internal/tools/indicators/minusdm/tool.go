// Package minusdm registers the minus_dm indicator with the talib dispatcher.
package minusdm

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "minus_dm",
		Description: "Minus Directional Movement",
		Group:       "momentum",
		Params:      talib.ParamsHLPeriod(14),
		Run:         talib.RunHLPeriod("minus_dm", 14, talib.MINUSDMFn),
	})
}
