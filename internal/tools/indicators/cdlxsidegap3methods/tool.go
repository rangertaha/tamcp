// Package cdlxsidegap3methods registers the cdlxsidegap3methods indicator with the talib dispatcher.
package cdlxsidegap3methods

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlxsidegap3methods",
		Description: "Up/Down-side Gap Three Methods candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlxsidegap3methods", talib.CDLXSIDEGAP3METHODSFn),
	})
}
