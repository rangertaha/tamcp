// Package sub registers the sub indicator with the talib dispatcher.
package sub

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sub",
		Description: "Vector Sub (a - b)",
		Group:       "operator",
		Params:      talib.ParamsTwoReal(),
		Run:         talib.RunTwoReal("sub", talib.SUBFn),
	})
}
