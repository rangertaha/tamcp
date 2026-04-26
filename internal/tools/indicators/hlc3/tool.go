// Package hlc3 registers the Typical Price helper (Pandas TA).
package hlc3

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "hlc3",
		Description: "Typical Price: (high + low + close) / 3.",
		Group:       "price",
		Params:      talib.ParamsHLC(),
		Run:         talib.RunHLC("hlc3", talib.HLC3Fn),
	})
}
