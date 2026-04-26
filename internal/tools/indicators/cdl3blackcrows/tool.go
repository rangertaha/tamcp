// Package cdl3blackcrows registers the cdl3blackcrows indicator with the talib dispatcher.
package cdl3blackcrows

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdl3blackcrows",
		Description: "Three Black Crows candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdl3blackcrows", talib.CDL3BLACKCROWSFn),
	})
}
