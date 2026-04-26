// Package bias registers the Bias indicator (Pandas TA).
package bias

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "bias",
		Description: "Bias: 100 * (close - SMA(close, p)) / SMA(close, p).",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(26),
		Run:         talib.RunRealPeriod("bias", 26, talib.BIASFn),
	})
}
