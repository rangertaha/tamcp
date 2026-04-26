// Package beta registers the beta indicator with the talib dispatcher.
package beta

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "beta",
		Description: "Beta",
		Group:       "statistic",
		Params:      talib.ParamsTwoRealPeriod(5),
		Run:         talib.RunTwoRealPeriod("beta", 5, talib.BETAFn),
	})
}
