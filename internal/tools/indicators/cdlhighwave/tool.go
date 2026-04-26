// Package cdlhighwave registers the cdlhighwave indicator with the talib dispatcher.
package cdlhighwave

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlhighwave",
		Description: "High-Wave candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlhighwave", talib.CDLHIGHWAVEFn),
	})
}
