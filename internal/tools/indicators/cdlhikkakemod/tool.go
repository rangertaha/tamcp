// Package cdlhikkakemod registers the cdlhikkakemod indicator with the talib dispatcher.
package cdlhikkakemod

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlhikkakemod",
		Description: "Hikkake Modified candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlhikkakemod", talib.CDLHIKKAKEMODFn),
	})
}
