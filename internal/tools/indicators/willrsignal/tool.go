// Package willrsignal registers Williams %R with EMA signal (utility).
package willrsignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "willr_signal",
		Description: "Williams %R with an EMA-smoothed signal line.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 14},
			{Name: "signal_period", Type: "int", Default: 9},
		},
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			w, sig := talib.WILLRSIGNALFn(h, l, c,
				talib.ArgInt(args, "period", 14),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(w, sig, [2]string{"willr", "signal"}), talib.Tersum("willr_signal", w), nil
		},
	})
}
