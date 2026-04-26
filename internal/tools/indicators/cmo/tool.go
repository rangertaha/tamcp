// Package cmo registers the cmo indicator with the talib dispatcher.
package cmo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cmo",
		Description: "Chande Momentum Oscillator",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("cmo", 14, talib.CMOFn),
	})
}
