// Package tan registers the tan indicator with the talib dispatcher.
package tan

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tan",
		Description: "Vector Trigonometric Tan",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("tan", talib.TANFn),
	})
}
