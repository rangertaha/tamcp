// Package cdldoji registers the cdldoji indicator with the talib dispatcher.
package cdldoji

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdldoji",
		Description: "Doji candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdldoji", talib.CDLDOJIFn),
	})
}
