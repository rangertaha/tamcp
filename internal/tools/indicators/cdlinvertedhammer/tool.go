// Package cdlinvertedhammer registers the cdlinvertedhammer indicator with the talib dispatcher.
package cdlinvertedhammer

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlinvertedhammer",
		Description: "Inverted Hammer candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlinvertedhammer", talib.CDLINVERTEDHAMMERFn),
	})
}
