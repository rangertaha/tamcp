// Package typprice registers the typprice indicator with the talib dispatcher.
package typprice

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "typprice",
		Description: "(high+low+close)/3",
		Group:       "price",
		Params:      talib.ParamsHLC(),
		Run:         talib.RunHLC("typprice", talib.TYPPRICEFn),
	})
}
