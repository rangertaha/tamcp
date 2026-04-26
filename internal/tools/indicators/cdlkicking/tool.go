// Package cdlkicking registers the cdlkicking indicator with the talib dispatcher.
package cdlkicking

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlkicking",
		Description: "Kicking candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlkicking", talib.CDLKICKINGFn),
	})
}
