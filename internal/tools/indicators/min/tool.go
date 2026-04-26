// Package min registers the min indicator with the talib dispatcher.
package min

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "min",
		Description: "Rolling minimum",
		Group:       "operator",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("min", 30, talib.MINFn),
	})
}
