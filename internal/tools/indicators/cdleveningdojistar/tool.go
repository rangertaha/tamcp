// Package cdleveningdojistar registers the cdleveningdojistar indicator with the talib dispatcher.
package cdleveningdojistar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdleveningdojistar",
		Description: "Evening Doji Star candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdleveningdojistar", talib.CDLEVENINGDOJISTARFn),
	})
}
