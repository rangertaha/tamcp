// Package medsmooth registers EMA-smoothed Median Price (utility).
package medsmooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "medprice_smoothed",
		Description: "EMA-smoothed Median Price (h+l)/2 over `period` bars.",
		Group:       "price",
		Params:      talib.ParamsHLPeriod(14),
		Run:         talib.RunHLPeriod("medprice_smoothed", 14, talib.MEDPRICESMOOTHFn),
	})
}
