// Package decycler registers Ehlers' High-Pass Decycler.
package decycler

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "decycler",
		Description: "Ehlers High-Pass Decycler: low-pass component (real - HP). Removes cycles shorter than `period`.",
		Group:       "cycle",
		Params:      talib.ParamsRealPeriod(50),
		Run:         talib.RunRealPeriod("decycler", 50, talib.DECYCLERFn),
	})
}
