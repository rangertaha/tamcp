// Package cdlengulfing registers the cdlengulfing indicator with the talib dispatcher.
package cdlengulfing

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlengulfing",
		Description: "Engulfing candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlengulfing", talib.CDLENGULFINGFn),
	})
}
