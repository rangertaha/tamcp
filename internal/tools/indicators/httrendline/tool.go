// Package httrendline registers the ht_trendline indicator with the talib dispatcher.
package httrendline

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ht_trendline",
		Description: "Hilbert Transform Instantaneous Trendline",
		Group:       "overlap",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("ht_trendline", talib.HTTRENDLINEFn),
	})
}
