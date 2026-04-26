// Package pdist registers the Price Distance indicator (Pandas TA).
package pdist

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pdist",
		Description: "Price Distance: 2*(high - low) - |close - open| + |open - close[-1]|.",
		Group:       "volatility",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunOHLC("pdist", talib.PDISTFn),
	})
}
