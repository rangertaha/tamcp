// Package asin registers the asin indicator with the talib dispatcher.
package asin

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "asin",
		Description: "Vector Trigonometric ASin",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("asin", talib.ASINFn),
	})
}
