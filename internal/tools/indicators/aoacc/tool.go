// Package aoacc registers Bill Williams' Acceleration/Deceleration (Pandas TA).
package aoacc

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ao_acc",
		Description: "Bill Williams' Acceleration/Deceleration: AO - SMA(AO, 5).",
		Group:       "momentum",
		Params:      talib.ParamsHL(),
		Run:         talib.RunHL("ao_acc", talib.AOACCFn),
	})
}
