// Package htdcphase registers the ht_dcphase indicator with the talib dispatcher.
package htdcphase

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ht_dcphase",
		Description: "Hilbert Transform Dominant Cycle Phase",
		Group:       "cycle",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("ht_dcphase", talib.HTDCPHASEFn),
	})
}
