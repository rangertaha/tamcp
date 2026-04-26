// Package retzscore registers Z-score of one-bar percent returns (utility).
package retzscore

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "return_zscore",
		Description: "Rolling Z-score of one-bar percent returns.",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("return_zscore", 30, talib.RETURNZSCOREFn),
	})
}
