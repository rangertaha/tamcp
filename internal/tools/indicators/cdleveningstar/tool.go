// Package cdleveningstar registers the cdleveningstar indicator with the talib dispatcher.
package cdleveningstar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdleveningstar",
		Description: "Evening Star candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdleveningstar", talib.CDLEVENINGSTARFn),
	})
}
