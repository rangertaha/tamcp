// Package wcpsmooth registers EMA-smoothed Weighted Close Price (utility).
package wcpsmooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "wcp_smoothed",
		Description: "EMA-smoothed Weighted Close Price (h+l+2c)/4.",
		Group:       "price",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("wcp_smoothed", 14, talib.WCPSMOOTHFn),
	})
}
