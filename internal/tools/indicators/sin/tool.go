// Package sin registers the sin indicator with the talib dispatcher.
package sin

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sin",
		Description: "Vector Trigonometric Sin",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("sin", talib.SINFn),
	})
}
