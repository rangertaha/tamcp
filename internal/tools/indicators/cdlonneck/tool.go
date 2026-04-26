// Package cdlonneck registers the cdlonneck indicator with the talib dispatcher.
package cdlonneck

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlonneck",
		Description: "On-Neck candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlonneck", talib.CDLONNECKFn),
	})
}
