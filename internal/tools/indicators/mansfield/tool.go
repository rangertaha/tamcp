// Package mansfield registers Mansfield Relative Strength.
package mansfield

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mansfield",
		Description: "Mansfield Relative Strength: 100 * (a/b - SMA(a/b, p)) / SMA(a/b, p). Compares asset `a` to benchmark `b`.",
		Group:       "momentum",
		Params:      talib.ParamsTwoRealPeriod(52),
		Run:         talib.RunTwoRealPeriod("mansfield", 52, talib.MANSFIELDFn),
	})
}
