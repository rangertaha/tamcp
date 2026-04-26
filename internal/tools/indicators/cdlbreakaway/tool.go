// Package cdlbreakaway registers the cdlbreakaway indicator with the talib dispatcher.
package cdlbreakaway

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlbreakaway",
		Description: "Breakaway candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlbreakaway", talib.CDLBREAKAWAYFn),
	})
}
