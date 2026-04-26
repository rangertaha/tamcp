// Package ceil registers the ceil indicator with the talib dispatcher.
package ceil

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ceil",
		Description: "Vector Ceil",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("ceil", talib.CEILFn),
	})
}
