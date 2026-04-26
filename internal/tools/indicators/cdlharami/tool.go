// Package cdlharami registers the cdlharami indicator with the talib dispatcher.
package cdlharami

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlharami",
		Description: "Harami candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlharami", talib.CDLHARAMIFn),
	})
}
