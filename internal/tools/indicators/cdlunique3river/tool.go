// Package cdlunique3river registers the cdlunique3river indicator with the talib dispatcher.
package cdlunique3river

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlunique3river",
		Description: "Unique Three River candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlunique3river", talib.CDLUNIQUE3RIVERFn),
	})
}
