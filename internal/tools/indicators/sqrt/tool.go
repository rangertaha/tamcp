// Package sqrt registers the sqrt indicator with the talib dispatcher.
package sqrt

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sqrt",
		Description: "Vector Square Root",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("sqrt", talib.SQRTFn),
	})
}
