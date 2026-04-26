// Package medprice registers the medprice indicator with the talib dispatcher.
package medprice

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "medprice",
		Description: "(high+low)/2",
		Group:       "price",
		Params:      talib.ParamsHL(),
		Run:         talib.RunHL("medprice", talib.MEDPRICEFn),
	})
}
