// Package mfv registers per-bar Money Flow Volume (Pandas TA).
package mfv

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mfv",
		Description: "Per-bar Money Flow Volume: ((c-l) - (h-c)) / (h-l) * volume.",
		Group:       "volume",
		Params:      talib.ParamsHLCV(),
		Run:         talib.RunHLCV("mfv", talib.MFVFn),
	})
}
