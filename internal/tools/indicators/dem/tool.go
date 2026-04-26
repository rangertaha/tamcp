// Package dem registers the DeMarker indicator (Pandas TA).
package dem

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "dem",
		Description: "DeMarker: SMA(demax, p) / (SMA(demax, p) + SMA(demin, p)) where demax = max(Δhigh, 0), demin = max(-Δlow, 0).",
		Group:       "momentum",
		Params:      talib.ParamsHLPeriod(14),
		Run:         talib.RunHLPeriod("dem", 14, talib.DEMFn),
	})
}
