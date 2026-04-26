// Package correl registers the correl indicator with the talib dispatcher.
package correl

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "correl",
		Description: "Pearson's correlation",
		Group:       "statistic",
		Params:      talib.ParamsTwoRealPeriod(30),
		Run:         talib.RunTwoRealPeriod("correl", 30, talib.CORRELFn),
	})
}
