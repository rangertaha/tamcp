// Package cdlbelthold registers the cdlbelthold indicator with the talib dispatcher.
package cdlbelthold

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlbelthold",
		Description: "Belt-hold candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlbelthold", talib.CDLBELTHOLDFn),
	})
}
