// Package chop registers the Choppiness Index (Pandas TA).
package chop

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "chop",
		Description: "Choppiness Index: 100 * log10(SUM(TR,p) / (HHV(high,p) - LLV(low,p))) / log10(p).",
		Group:       "volatility",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("chop", 14, talib.CHOPFn),
	})
}
