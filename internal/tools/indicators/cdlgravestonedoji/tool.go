// Package cdlgravestonedoji registers the cdlgravestonedoji indicator with the talib dispatcher.
package cdlgravestonedoji

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlgravestonedoji",
		Description: "Gravestone Doji candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlgravestonedoji", talib.CDLGRAVESTONEDOJIFn),
	})
}
