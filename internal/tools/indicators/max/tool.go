// Package max registers the max indicator with the talib dispatcher.
package max

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "max",
		Description: "Rolling maximum",
		Group:       "operator",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("max", 30, talib.MAXFn),
	})
}
