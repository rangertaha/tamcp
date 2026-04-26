// Package trendflex registers Ehlers' TrendFlex indicator (Pandas TA).
package trendflex

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "trendflex",
		Description: "Ehlers TrendFlex: SSF-prefiltered slope normalised by its rolling RMS.",
		Group:       "momentum",
		Params:      talib.ParamsRealPeriod(20),
		Run:         talib.RunRealPeriod("trendflex", 20, talib.TRENDFLEXFn),
	})
}
