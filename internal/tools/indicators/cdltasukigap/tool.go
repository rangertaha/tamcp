// Package cdltasukigap registers the cdltasukigap indicator with the talib dispatcher.
package cdltasukigap

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdltasukigap",
		Description: "Tasuki Gap candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdltasukigap", talib.CDLTASUKIGAPFn),
	})
}
