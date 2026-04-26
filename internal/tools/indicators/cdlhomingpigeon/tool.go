// Package cdlhomingpigeon registers the cdlhomingpigeon indicator with the talib dispatcher.
package cdlhomingpigeon

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlhomingpigeon",
		Description: "Homing Pigeon candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlhomingpigeon", talib.CDLHOMINGPIGEONFn),
	})
}
