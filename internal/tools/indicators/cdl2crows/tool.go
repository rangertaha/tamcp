// Package cdl2crows registers the cdl2crows indicator with the talib dispatcher.
package cdl2crows

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdl2crows",
		Description: "Two Crows candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdl2crows", talib.CDL2CROWSFn),
	})
}
