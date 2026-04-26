// Package swma registers the Symmetric Weighted MA (Pandas TA).
package swma

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "swma",
		Description: "Symmetric Weighted MA: classic 4-bar [1,2,2,1]/6 filter.",
		Group:       "overlap",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("swma", talib.SWMAFn),
	})
}
