// Package atan registers the atan indicator with the talib dispatcher.
package atan

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "atan",
		Description: "Vector Trigonometric ATan",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("atan", talib.ATANFn),
	})
}
