// Package gaussian registers Ehlers' 4-pole Gaussian Filter.
package gaussian

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "gaussian",
		Description: "Ehlers 4-pole Gaussian Filter: sharper roll-off than EMA at the same nominal period.",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("gaussian", 20, talib.GAUSSIANFn),
	})
}
