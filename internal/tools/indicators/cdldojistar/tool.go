// Package cdldojistar registers the cdldojistar indicator with the talib dispatcher.
package cdldojistar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdldojistar",
		Description: "Doji Star candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdldojistar", talib.CDLDOJISTARFn),
	})
}
