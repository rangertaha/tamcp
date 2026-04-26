// Package floor registers the floor indicator with the talib dispatcher.
package floor

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "floor",
		Description: "Vector Floor",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("floor", talib.FLOORFn),
	})
}
