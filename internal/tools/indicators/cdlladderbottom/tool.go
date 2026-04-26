// Package cdlladderbottom registers the cdlladderbottom indicator with the talib dispatcher.
package cdlladderbottom

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlladderbottom",
		Description: "Ladder Bottom candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlladderbottom", talib.CDLLADDERBOTTOMFn),
	})
}
