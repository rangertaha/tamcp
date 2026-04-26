// Package cdlspinningtop registers the cdlspinningtop indicator with the talib dispatcher.
package cdlspinningtop

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlspinningtop",
		Description: "Spinning Top candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlspinningtop", talib.CDLSPINNINGTOPFn),
	})
}
