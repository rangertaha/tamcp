// Package cdlmorningstar registers the cdlmorningstar indicator with the talib dispatcher.
package cdlmorningstar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlmorningstar",
		Description: "Morning Star candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlmorningstar", talib.CDLMORNINGSTARFn),
	})
}
