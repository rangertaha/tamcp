// Package rocp registers the rocp indicator with the talib dispatcher.
package rocp

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rocp",
		Description: "Rate of change Percentage",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("rocp", 10, talib.ROCPFn),
	})
}
