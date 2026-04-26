// Package tpsmooth registers EMA-smoothed Typical Price (utility).
package tpsmooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tp_smoothed",
		Description: "EMA of Typical Price (HLC3).",
		Group:       "price",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("tp_smoothed", 14, talib.TPSMOOTHFn),
	})
}
