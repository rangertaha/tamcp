// Package adxr registers the adxr indicator with the talib dispatcher.
package adxr

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "adxr",
		Description: "Average Directional Movement Index Rating",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("adxr", 14, talib.ADXRFn),
	})
}
