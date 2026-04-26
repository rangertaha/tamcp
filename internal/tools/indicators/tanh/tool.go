// Package tanh registers the tanh indicator with the talib dispatcher.
package tanh

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "tanh",
		Description: "Vector Trigonometric Tanh",
		Group:       "math",
		Params:      talib.ParamsRealOnly(),
		Run:         talib.RunRealOnly("tanh", talib.TANHFn),
	})
}
