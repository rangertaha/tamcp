// Package cdllongleggeddoji registers the cdllongleggeddoji indicator with the talib dispatcher.
package cdllongleggeddoji

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdllongleggeddoji",
		Description: "Long-Legged Doji candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdllongleggeddoji", talib.CDLLONGLEGGEDDOJIFn),
	})
}
