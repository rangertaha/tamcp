// Package cdlabandonedbaby registers the cdlabandonedbaby indicator with the talib dispatcher.
package cdlabandonedbaby

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlabandonedbaby",
		Description: "Abandoned Baby candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlabandonedbaby", talib.CDLABANDONEDBABYFn),
	})
}
