// Package donchianpct registers a Donchian %B indicator.
package donchianpct

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "donchian_pct",
		Description: "Donchian %B: (close - lower) / (upper - lower) over `period` bars.",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(20),
		Run:         talib.RunHLCPeriod("donchian_pct", 20, talib.DONCHIANPCTFn),
	})
}
