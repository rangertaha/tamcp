// Package pgo registers the Pretty Good Oscillator (Pandas TA).
package pgo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pgo",
		Description: "Pretty Good Oscillator: (close - SMA(close,p)) / EMA(TR, p).",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("pgo", 14, talib.PGOFn),
	})
}
