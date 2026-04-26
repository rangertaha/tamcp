// Package cdllongline registers the cdllongline indicator with the talib dispatcher.
package cdllongline

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdllongline",
		Description: "Long Line candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdllongline", talib.CDLLONGLINEFn),
	})
}
