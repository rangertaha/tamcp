// Package vwappct registers % deviation of close from running VWAP.
package vwappct

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vwap_pct",
		Description: "Percentage deviation of close from the cumulative VWAP: 100 * (close - VWAP) / VWAP.",
		Group:       "volume",
		Params:      talib.ParamsHLCV(),
		Run:         talib.RunHLCV("vwap_pct", talib.VWAPPCTFn),
	})
}
