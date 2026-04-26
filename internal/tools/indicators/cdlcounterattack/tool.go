// Package cdlcounterattack registers the cdlcounterattack indicator with the talib dispatcher.
package cdlcounterattack

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlcounterattack",
		Description: "Counterattack candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlcounterattack", talib.CDLCOUNTERATTACKFn),
	})
}
