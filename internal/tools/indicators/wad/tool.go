// Package wad registers Williams Accumulation/Distribution (cinar).
package wad

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "wad",
		Description: "Williams Accumulation/Distribution: cumulative sum of price action vs the prior close's true high/low.",
		Group:       "volume",
		Params:      talib.ParamsHLC(),
		Run:         talib.RunHLC("wad", talib.WADFn),
	})
}
