// Package cdl3linestrike registers the cdl3linestrike indicator with the talib dispatcher.
package cdl3linestrike

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdl3linestrike",
		Description: "Three Line Strike candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdl3linestrike", talib.CDL3LINESTRIKEFn),
	})
}
