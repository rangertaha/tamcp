// Package cdlthrusting registers the cdlthrusting indicator with the talib dispatcher.
package cdlthrusting

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlthrusting",
		Description: "Thrusting candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlthrusting", talib.CDLTHRUSTINGFn),
	})
}
