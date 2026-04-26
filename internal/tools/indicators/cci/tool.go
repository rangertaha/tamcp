// Package cci registers the cci indicator with the talib dispatcher.
package cci

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cci",
		Description: "Commodity Channel Index",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("cci", 14, talib.CCIFn),
	})
}
