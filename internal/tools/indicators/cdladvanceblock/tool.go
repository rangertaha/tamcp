// Package cdladvanceblock registers the cdladvanceblock indicator with the talib dispatcher.
package cdladvanceblock

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdladvanceblock",
		Description: "Advance Block candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdladvanceblock", talib.CDLADVANCEBLOCKFn),
	})
}
