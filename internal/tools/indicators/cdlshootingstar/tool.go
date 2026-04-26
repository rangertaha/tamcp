// Package cdlshootingstar registers the cdlshootingstar indicator with the talib dispatcher.
package cdlshootingstar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlshootingstar",
		Description: "Shooting Star candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlshootingstar", talib.CDLSHOOTINGSTARFn),
	})
}
