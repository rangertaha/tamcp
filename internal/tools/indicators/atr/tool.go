// Package atr registers the atr indicator with the talib dispatcher.
package atr

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "atr",
		Description: "Average True Range",
		Group:       "volatility",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("atr", 14, talib.ATRFn),
	})
}
