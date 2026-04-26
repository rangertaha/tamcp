// Package decay registers the Linear Decay indicator (Pandas TA).
package decay

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "decay",
		Description: "Linear decay: out[i] = max(real[i], out[i-1] - 1/length, 0). Useful for fading triggers.",
		Group:       "trend",
		Params:      talib.ParamsRealPeriod(5),
		Run:         talib.RunRealPeriod("decay", 5, talib.DECAYFn),
	})
}
