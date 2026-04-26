// Package cdlrickshawman registers the cdlrickshawman indicator with the talib dispatcher.
package cdlrickshawman

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlrickshawman",
		Description: "Rickshaw Man candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlrickshawman", talib.CDLRICKSHAWMANFn),
	})
}
