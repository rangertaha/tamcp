// Package pwma registers the Pascal Weighted MA (Pandas TA).
package pwma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pwma",
		Description: "Pascal Weighted MA: weights are binomial coefficients C(p-1, k); central sample carries the heaviest weight.",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("pwma", 10, talib.PWMAFn),
	})
}
