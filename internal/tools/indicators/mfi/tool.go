// Package mfi registers the mfi indicator with the talib dispatcher.
package mfi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mfi",
		Description: "Money Flow Index",
		Group:       "momentum",
		Params:      talib.ParamsHLCVPeriod(14),
		Run:         talib.RunHLCVPeriod("mfi", 14, talib.MFIFn),
	})
}
