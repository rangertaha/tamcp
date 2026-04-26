// Package cdlhammer registers the cdlhammer indicator with the talib dispatcher.
package cdlhammer

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlhammer",
		Description: "Hammer candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlhammer", talib.CDLHAMMERFn),
	})
}
