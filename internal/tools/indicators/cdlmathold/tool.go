// Package cdlmathold registers the cdlmathold indicator with the talib dispatcher.
package cdlmathold

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlmathold",
		Description: "Mat Hold candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlmathold", talib.CDLMATHOLDFn),
	})
}
