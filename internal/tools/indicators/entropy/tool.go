// Package entropy registers the Shannon entropy indicator (Pandas TA).
package entropy

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "entropy",
		Description: "Shannon entropy (natural log) of the sliding window normalised to a probability distribution. Negative or zero values contribute nothing.",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("entropy", 10, talib.ENTROPYFn),
	})
}
