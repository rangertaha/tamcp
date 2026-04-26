// Package cdlhangingman registers the cdlhangingman indicator with the talib dispatcher.
package cdlhangingman

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlhangingman",
		Description: "Hanging Man candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlhangingman", talib.CDLHANGINGMANFn),
	})
}
