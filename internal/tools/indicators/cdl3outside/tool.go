// Package cdl3outside registers the cdl3outside indicator with the talib dispatcher.
package cdl3outside

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdl3outside",
		Description: "Three Outside Up/Down candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdl3outside", talib.CDL3OUTSIDEFn),
	})
}
