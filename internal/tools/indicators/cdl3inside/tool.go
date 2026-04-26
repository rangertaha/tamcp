// Package cdl3inside registers the cdl3inside indicator with the talib dispatcher.
package cdl3inside

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdl3inside",
		Description: "Three Inside Up/Down candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdl3inside", talib.CDL3INSIDEFn),
	})
}
