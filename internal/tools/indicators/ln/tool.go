// Package ln registers the ln indicator with the talib dispatcher.
package ln

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ln",
		Description: "Vector Log Natural",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("ln", talib.LNFn),
	})
}
