// Package cdldragonflydoji registers the cdldragonflydoji indicator with the talib dispatcher.
package cdldragonflydoji

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdldragonflydoji",
		Description: "Dragonfly Doji candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdldragonflydoji", talib.CDLDRAGONFLYDOJIFn),
	})
}
