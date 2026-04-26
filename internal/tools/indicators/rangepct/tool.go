// Package rangepct registers Range as a percentage of close (utility).
package rangepct

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "range_pct",
		Description: "Bar range as a percentage of close: 100 * (high - low) / close.",
		Group:       "volatility",
		Params:      talib.ParamsHLC(),
		Run:         talib.RunHLC("range_pct", talib.RANGEPCTFn),
	})
}
