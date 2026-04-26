// Package trangesmooth registers EMA-smoothed True Range (utility).
package trangesmooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "trange_smoothed",
		Description: "EMA-smoothed True Range over `period` bars.",
		Group:       "volatility",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("trange_smoothed", 14, talib.TRANGESMOOTHFn),
	})
}
