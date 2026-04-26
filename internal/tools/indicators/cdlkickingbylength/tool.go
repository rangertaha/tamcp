// Package cdlkickingbylength registers the cdlkickingbylength indicator with the talib dispatcher.
package cdlkickingbylength

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlkickingbylength",
		Description: "Kicking by Length candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlkickingbylength", talib.CDLKICKINGBYLENGTHFn),
	})
}
