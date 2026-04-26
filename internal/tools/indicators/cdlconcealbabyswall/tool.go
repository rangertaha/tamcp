// Package cdlconcealbabyswall registers the cdlconcealbabyswall indicator with the talib dispatcher.
package cdlconcealbabyswall

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlconcealbabyswall",
		Description: "Concealing Baby Swallow candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlconcealbabyswall", talib.CDLCONCEALBABYSWALLFn),
	})
}
