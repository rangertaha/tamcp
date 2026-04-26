// Package ao registers Bill Williams' Awesome Oscillator (Pandas TA, cinar).
package ao

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ao",
		Description: "Awesome Oscillator: SMA((H+L)/2, 5) − SMA((H+L)/2, 34).",
		Group:       "momentum",
		Params:      talib.ParamsHL(),
		Run:         talib.RunHL("ao", talib.AOFn),
	})
}
