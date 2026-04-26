// Package cdlharamicross registers the cdlharamicross indicator with the talib dispatcher.
package cdlharamicross

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlharamicross",
		Description: "Harami Cross candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlharamicross", talib.CDLHARAMICROSSFn),
	})
}
