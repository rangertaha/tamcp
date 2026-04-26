// Package tmf registers the Twiggs Money Flow indicator (cinar).
package tmf

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tmf",
		Description: "Twiggs Money Flow: EMA(volume * money_flow_multiplier with prior-close trueHigh/Low) / EMA(volume).",
		Group:       "volume",
		Params:      talib.ParamsHLCVPeriod(21),
		Run:         talib.RunHLCVPeriod("tmf", 21, talib.TMFFn),
	})
}
