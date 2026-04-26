// Package mom registers the mom indicator with the talib dispatcher.
package mom

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mom",
		Description: "Momentum",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("mom", 10, talib.MOMFn),
	})
}
