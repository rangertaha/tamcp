// Package adzscore registers Z-score of the AD line (utility).
package adzscore

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ad_zscore",
		Description: "Rolling Z-score of the Accumulation/Distribution line.",
		Group:       "volume",
		Params:      talib.ParamsHLCVPeriod(30),
		Run:         talib.RunHLCVPeriod("ad_zscore", 30, talib.ADZSCOREFn),
	})
}
