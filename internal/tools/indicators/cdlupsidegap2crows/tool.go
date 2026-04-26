// Package cdlupsidegap2crows registers the cdlupsidegap2crows indicator with the talib dispatcher.
package cdlupsidegap2crows

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlupsidegap2crows",
		Description: "Upside Gap Two Crows candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlupsidegap2crows", talib.CDLUPSIDEGAP2CROWSFn),
	})
}
