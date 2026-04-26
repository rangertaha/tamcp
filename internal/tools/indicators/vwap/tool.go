// Package vwap registers the cumulative VWAP indicator (Pandas TA, cinar).
package vwap

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vwap",
		Description: "Cumulative Volume Weighted Average Price using HLC3 typical price. No session reset; pre-slice inputs per session if needed.",
		Group:       "volume",
		Params:      talib.ParamsHLCV(),
		Run:         talib.RunHLCV("vwap", talib.VWAPFn),
	})
}
