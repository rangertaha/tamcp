// Package htdcperiod registers the ht_dcperiod indicator with the talib dispatcher.
package htdcperiod

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ht_dcperiod",
		Description: "Hilbert Transform Dominant Cycle Period",
		Group:       "cycle",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("ht_dcperiod", talib.HTDCPERIODFn),
	})
}
