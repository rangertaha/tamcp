// Package plusdi registers the plus_di indicator with the talib dispatcher.
package plusdi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "plus_di",
		Description: "Plus Directional Indicator",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("plus_di", 14, talib.PLUSDIFn),
	})
}
