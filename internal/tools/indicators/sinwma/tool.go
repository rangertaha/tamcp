// Package sinwma registers the Sine Weighted MA (Pandas TA).
package sinwma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sinwma",
		Description: "Sine Weighted MA: weights w[k] = sin((k+1)·π / (p+1)), centred bell-shape.",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(14),
		Run:         talib.RunRealPeriod("sinwma", 14, talib.SINWMAFn),
	})
}
