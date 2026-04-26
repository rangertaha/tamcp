// Package cdlseparatinglines registers the cdlseparatinglines indicator with the talib dispatcher.
package cdlseparatinglines

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlseparatinglines",
		Description: "Separating Lines candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlseparatinglines", talib.CDLSEPARATINGLINESFn),
	})
}
