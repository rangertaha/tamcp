// Package cdlsticksandwich registers the cdlsticksandwich indicator with the talib dispatcher.
package cdlsticksandwich

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlsticksandwich",
		Description: "Stick Sandwich candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlsticksandwich", talib.CDLSTICKSANDWICHFn),
	})
}
