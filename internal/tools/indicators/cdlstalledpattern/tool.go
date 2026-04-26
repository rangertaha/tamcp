// Package cdlstalledpattern registers the cdlstalledpattern indicator with the talib dispatcher.
package cdlstalledpattern

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlstalledpattern",
		Description: "Stalled Pattern candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlstalledpattern", talib.CDLSTALLEDPATTERNFn),
	})
}
