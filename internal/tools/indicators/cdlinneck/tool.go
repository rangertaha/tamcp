// Package cdlinneck registers the cdlinneck indicator with the talib dispatcher.
package cdlinneck

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlinneck",
		Description: "In-Neck candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlinneck", talib.CDLINNECKFn),
	})
}
