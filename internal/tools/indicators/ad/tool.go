// Package ad registers the ad indicator with the talib dispatcher.
package ad

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ad",
		Description: "Chaikin A/D Line",
		Group:       "volume",
		Params:      talib.ParamsHLCV(),
		Run:         talib.RunHLCV("ad", talib.ADFn),
	})
}
