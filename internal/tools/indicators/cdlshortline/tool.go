// Package cdlshortline registers the cdlshortline indicator with the talib dispatcher.
package cdlshortline

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlshortline",
		Description: "Short Line candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlshortline", talib.CDLSHORTLINEFn),
	})
}
