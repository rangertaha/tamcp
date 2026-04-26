// Package cfo registers the Chande Forecast Oscillator (Pandas TA).
package cfo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cfo",
		Description: "Chande Forecast Oscillator: 100 * (real - TSF(real, p)) / real.",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(9),
		Run:         talib.RunRealPeriod("cfo", 9, talib.CFOFn),
	})
}
