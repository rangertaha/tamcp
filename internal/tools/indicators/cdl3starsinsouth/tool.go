// Package cdl3starsinsouth registers the cdl3starsinsouth indicator with the talib dispatcher.
package cdl3starsinsouth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdl3starsinsouth",
		Description: "Three Stars in the South candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdl3starsinsouth", talib.CDL3STARSINSOUTHFn),
	})
}
