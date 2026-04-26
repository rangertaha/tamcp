// Package wclprice registers the wclprice indicator with the talib dispatcher.
package wclprice

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "wclprice",
		Description: "(high+low+2*close)/4",
		Group:       "price",
		Params:      talib.ParamsHLC(),
		Run:         talib.RunHLC("wclprice", talib.WCLPRICEFn),
	})
}
