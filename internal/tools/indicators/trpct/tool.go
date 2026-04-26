// Package trpct registers True Range as a percentage of close (utility).
package trpct

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "true_range_pct",
		Description: "True Range as a percentage of close: 100 * TR / close.",
		Group:       "volatility",
		Params:      talib.ParamsHLC(),
		Run:         talib.RunHLC("true_range_pct", talib.TRUERANGEPCTFn),
	})
}
