// Package adx registers the adx indicator with the talib dispatcher.
package adx

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "adx",
		Description: "Average Directional Movement Index",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("adx", 14, talib.ADXFn),
	})
}
