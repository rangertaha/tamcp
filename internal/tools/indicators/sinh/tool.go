// Package sinh registers the sinh indicator with the talib dispatcher.
package sinh

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sinh",
		Description: "Vector Trigonometric Sinh",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("sinh", talib.SINHFn),
	})
}
