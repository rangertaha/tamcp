// Package natr registers the natr indicator with the talib dispatcher.
package natr

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "natr",
		Description: "Normalized Average True Range",
		Group:       "volatility",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("natr", 14, talib.NATRFn),
	})
}
