// Package cdlpiercing registers the cdlpiercing indicator with the talib dispatcher.
package cdlpiercing

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlpiercing",
		Description: "Piercing candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlpiercing", talib.CDLPIERCINGFn),
	})
}
