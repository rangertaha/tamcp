// Package add registers the add indicator with the talib dispatcher.
package add

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "add",
		Description: "Vector Add (a + b)",
		Group:       "operator",
		Params:      talib.ParamsTwoReal(),
		Run:         talib.RunTwoReal("add", talib.ADDFn),
	})
}
