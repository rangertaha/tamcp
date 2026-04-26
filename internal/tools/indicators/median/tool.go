// Package median registers the rolling median indicator (Pandas TA).
package median

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "median",
		Description: "Rolling median over `period` bars.",
		Group:       "statistic",
		Params:      talib.ParamsRealPeriod(30),
		Run:         talib.RunRealPeriod("median", 30, talib.MEDIANFn),
	})
}
