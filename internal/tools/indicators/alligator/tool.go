// Package alligator registers Bill Williams' Alligator (Pandas TA).
package alligator

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "alligator",
		Description: "Bill Williams Alligator: jaw=SMMA(13) lagged 8, teeth=SMMA(8) lagged 5, lips=SMMA(5) lagged 3 of (H+L)/2.",
		Group:       "overlap",
		Params:      talib.ParamsHL(),
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			j, t, lp := talib.ALLIGATORFn(h, l)
			return talib.Three(j, t, lp, [3]string{"jaw", "teeth", "lips"}), talib.Tersum("alligator", t), nil
		},
	})
}
