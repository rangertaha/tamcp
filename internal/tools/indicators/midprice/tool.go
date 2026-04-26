// Package midprice registers the midprice indicator with the talib dispatcher.
package midprice

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "midprice",
		Description: "MidPoint Price over period (high/low)",
		Group:       "overlap",
		Params:      talib.ParamsHLPeriod(14),
		Run:         talib.RunHLPeriod("midprice", 14, talib.MIDPRICEFn),
	})
}
