// Package httrendmode registers the ht_trendmode indicator with the talib dispatcher.
package httrendmode

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ht_trendmode",
		Description: "Hilbert Transform trend/cycle mode",
		Group:       "cycle",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("ht_trendmode", talib.HTTRENDMODEFn),
	})
}
