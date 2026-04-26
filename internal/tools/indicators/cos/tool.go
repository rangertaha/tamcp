// Package cos registers the cos indicator with the talib dispatcher.
package cos

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cos",
		Description: "Vector Trigonometric Cos",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("cos", talib.COSFn),
	})
}
