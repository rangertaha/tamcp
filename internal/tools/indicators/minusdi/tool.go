// Package minusdi registers the minus_di indicator with the talib dispatcher.
package minusdi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "minus_di",
		Description: "Minus Directional Indicator",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("minus_di", 14, talib.MINUSDIFn),
	})
}
