// Package avgprice registers the avgprice indicator with the talib dispatcher.
package avgprice

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "avgprice",
		Description: "(open+high+low+close)/4",
		Group:       "price",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunOHLC("avgprice", talib.AVGPRICEFn),
	})
}
