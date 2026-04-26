// Package ratrpct registers TR / ATR ratio (utility).
package ratrpct

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "range_atr_pct",
		Description: "TR / ATR(h,l,c, period). >1 means today's true range is larger than the recent average.",
		Group:       "volatility",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("range_atr_pct", 14, talib.RANGEATRPCTFn),
	})
}
