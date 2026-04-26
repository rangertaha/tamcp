// Package slopepct registers LINREG slope as a percentage of its value (utility).
package slopepct

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "slope_pct",
		Description: "100 * LINREG slope / LINREG value over `period` bars. Comparable across instruments.",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("slope_pct", 20, talib.SLOPEPCTFn),
	})
}
