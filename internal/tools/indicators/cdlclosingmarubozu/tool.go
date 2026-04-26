// Package cdlclosingmarubozu registers the cdlclosingmarubozu indicator with the talib dispatcher.
package cdlclosingmarubozu

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlclosingmarubozu",
		Description: "Closing Marubozu candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlclosingmarubozu", talib.CDLCLOSINGMARUBOZUFn),
	})
}
