// Package cdlidentical3crows registers the cdlidentical3crows indicator with the talib dispatcher.
package cdlidentical3crows

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlidentical3crows",
		Description: "Identical Three Crows candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlidentical3crows", talib.CDLIDENTICAL3CROWSFn),
	})
}
