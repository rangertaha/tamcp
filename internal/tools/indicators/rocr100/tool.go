// Package rocr100 registers the rocr100 indicator with the talib dispatcher.
package rocr100

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rocr100",
		Description: "Rate of change Ratio (×100)",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("rocr100", 10, talib.ROCR100Fn),
	})
}
