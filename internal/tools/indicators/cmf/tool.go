// Package cmf registers the Chaikin Money Flow indicator (Pandas TA, cinar).
package cmf

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cmf",
		Description: "Chaikin Money Flow: SUM(money_flow_volume, p) / SUM(volume, p), where money_flow_volume = ((c-l)-(h-c))/(h-l) * volume.",
		Group:       "volume",
		Params:      talib.ParamsHLCVPeriod(20),
		Run:         talib.RunHLCVPeriod("cmf", 20, talib.CMFFn),
	})
}
