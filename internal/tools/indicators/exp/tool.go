// Package exp registers the exp indicator with the talib dispatcher.
package exp

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "exp",
		Description: "Vector Arithmetic Exp",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("exp", talib.EXPFn),
	})
}
