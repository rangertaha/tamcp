// Package trange registers the trange indicator with the talib dispatcher.
package trange

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "trange",
		Description: "True Range",
		Group:       "volatility",
		Params:      talib.ParamsHLC(),
		Run:         talib.RunHLC("trange", talib.TRANGEFn),
	})
}
