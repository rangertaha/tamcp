// Package vhf registers the Vertical Horizontal Filter (Pandas TA).
package vhf

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vhf",
		Description: "Vertical Horizontal Filter: (HHV(close,p) - LLV(close,p)) / SUM(|Δclose|, p).",
		Group:       "trend",
		Params:      talib.ParamsRealPeriod(28),
		Run:         talib.RunRealPeriod("vhf", 28, talib.VHFFn),
	})
}
