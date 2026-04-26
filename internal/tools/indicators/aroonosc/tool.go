// Package aroonosc registers the aroonosc indicator with the talib dispatcher.
package aroonosc

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "aroonosc",
		Description: "Aroon Oscillator",
		Group:       "momentum",
		Params:      talib.ParamsHLPeriod(14),
		Run:         talib.RunHLPeriod("aroonosc", 14, talib.AROONOSCFn),
	})
}
