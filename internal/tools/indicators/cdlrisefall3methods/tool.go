// Package cdlrisefall3methods registers the cdlrisefall3methods indicator with the talib dispatcher.
package cdlrisefall3methods

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlrisefall3methods",
		Description: "Rising/Falling Three Methods candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlrisefall3methods", talib.CDLRISEFALL3METHODSFn),
	})
}
