// Package acos registers the acos indicator with the talib dispatcher.
package acos

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "acos",
		Description: "Vector Trigonometric ACos",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("acos", talib.ACOSFn),
	})
}
