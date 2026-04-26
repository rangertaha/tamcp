// Package log10 registers the log10 indicator with the talib dispatcher.
package log10

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "log10",
		Description: "Vector Log10",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("log10", talib.LOG10Fn),
	})
}
