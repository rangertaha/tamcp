// Package cdltristar registers the cdltristar indicator with the talib dispatcher.
package cdltristar

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdltristar",
		Description: "Tristar candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdltristar", talib.CDLTRISTARFn),
	})
}
