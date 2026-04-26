// Package cdlmatchinglow registers the cdlmatchinglow indicator with the talib dispatcher.
package cdlmatchinglow

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdlmatchinglow",
		Description: "Matching Low candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdlmatchinglow", talib.CDLMATCHINGLOWFn),
	})
}
