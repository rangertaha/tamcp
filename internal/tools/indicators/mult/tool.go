// Package mult registers the mult indicator with the talib dispatcher.
package mult

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mult",
		Description: "Vector Mult (a * b)",
		Group:       "operator",
		Params:      talib.ParamsTwoReal(),
		Run:         talib.RunTwoReal("mult", talib.MULTFn),
	})
}
