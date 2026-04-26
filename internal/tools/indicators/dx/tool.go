// Package dx registers the dx indicator with the talib dispatcher.
package dx

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "dx",
		Description: "Directional Movement Index",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
		Run:         talib.RunHLCPeriod("dx", 14, talib.DXFn),
	})
}
