// Package cosh registers the cosh indicator with the talib dispatcher.
package cosh

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cosh",
		Description: "Vector Trigonometric Cosh",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("cosh", talib.COSHFn),
	})
}
