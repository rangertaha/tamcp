// Package cdltakuri registers the cdltakuri indicator with the talib dispatcher.
package cdltakuri

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdltakuri",
		Description: "Takuri candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdltakuri", talib.CDLTAKURIFn),
	})
}
