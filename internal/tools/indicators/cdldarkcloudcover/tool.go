// Package cdldarkcloudcover registers the cdldarkcloudcover indicator with the talib dispatcher.
package cdldarkcloudcover

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdldarkcloudcover",
		Description: "Dark Cloud Cover candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdldarkcloudcover", talib.CDLDARKCLOUDCOVERFn),
	})
}
