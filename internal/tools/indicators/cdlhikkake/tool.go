// Package cdlhikkake registers the cdlhikkake indicator with the talib dispatcher.
package cdlhikkake

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlhikkake",
		Description: "Hikkake candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlhikkake", talib.CDLHIKKAKEFn),
	})
}
