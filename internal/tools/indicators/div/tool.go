// Package div registers the div indicator with the talib dispatcher.
package div

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "div",
		Description: "Vector Div (a / b)",
		Group:       "operator",
		Params:      talib.ParamsTwoReal(),
		Run:         talib.RunTwoReal("div", talib.DIVFn),
	})
}
