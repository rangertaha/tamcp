// Package vortexdiff registers VI+ - VI- (Vortex spread).
package vortexdiff

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vortex_diff",
		Description: "Vortex spread: VI+ - VI- over `period` bars.",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("vortex_diff", 14, talib.VORTEXDIFFFn),
	})
}
