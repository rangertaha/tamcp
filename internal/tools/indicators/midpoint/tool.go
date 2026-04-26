// Package midpoint registers the midpoint indicator with the talib dispatcher.
package midpoint

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "midpoint",
		Description: "MidPoint over period: (max+min)/2",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("midpoint", 14, talib.MIDPOINTFn),
	})
}
