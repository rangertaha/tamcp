// Package bop registers the bop indicator with the talib dispatcher.
package bop

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "bop",
		Description: "Balance of Power",
		Group:       "momentum",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunOHLC("bop", talib.BOPFn),
	})
}
