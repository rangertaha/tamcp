// Package ohlc4 registers the OHLC Average price helper (Pandas TA).
package ohlc4

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ohlc4",
		Description: "OHLC Average: (open + high + low + close) / 4.",
		Group:       "price",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunOHLC("ohlc4", talib.OHLC4Fn),
	})
}
