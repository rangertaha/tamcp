// Package cdlmarubozu registers the cdlmarubozu indicator with the talib dispatcher.
package cdlmarubozu

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlmarubozu",
		Description: "Marubozu candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlmarubozu", talib.CDLMARUBOZUFn),
	})
}
