// Package cdlgapsidesidewhite registers the cdlgapsidesidewhite indicator with the talib dispatcher.
package cdlgapsidesidewhite

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlgapsidesidewhite",
		Description: "Gap Side-by-Side White candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlgapsidesidewhite", talib.CDLGAPSIDESIDEWHITEFn),
	})
}
