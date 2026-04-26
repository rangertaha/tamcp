// Package cdlmorningdojistar registers the cdlmorningdojistar indicator with the talib dispatcher.
package cdlmorningdojistar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlmorningdojistar",
		Description: "Morning Doji Star candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlmorningdojistar", talib.CDLMORNINGDOJISTARFn),
	})
}
