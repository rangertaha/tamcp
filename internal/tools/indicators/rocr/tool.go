// Package rocr registers the rocr indicator with the talib dispatcher.
package rocr

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rocr",
		Description: "Rate of change Ratio",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("rocr", 10, talib.ROCRFn),
	})
}
